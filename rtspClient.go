package goRtspClient

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ardroh/goRtspClient/client"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/data"
	"github.com/ardroh/goRtspClient/handlers"
	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/readers"
	"github.com/ardroh/goRtspClient/responses"
)

type RtspClient struct {
	connection         net.Conn
	params             client.RtspConnectionParams
	rtspConnReader     *readers.RtspSocketReader
	rtspResponseBuffer responses.RtspResponseBuffer
	context            client.RtspConnectionContext
	DataPacketChan     chan data.DataPacket
}

func (client *RtspClient) Connect(connectionParams client.RtspConnectionParams) error {
	if client.connection != nil {
		return errors.New("Already connected")
	}
	client.params = connectionParams
	dialer := net.Dialer{Timeout: time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", connectionParams.IP, connectionParams.Port))
	if err != nil {
		return err
	}
	client.DataPacketChan = make(chan data.DataPacket)
	client.connection = conn
	client.rtspConnReader = readers.CreateRtspConnReader(client.connection)
	go client.rtspConnReader.StartReading()
	request := &handlers.RtspConnectRequest{
		RtspClient:   client,
		Transmission: connectionParams.Transmission,
		Transport:    connectionParams.Transport,
	}
	handler := handlers.RtspHandlerFactory{}.Create()
	handler.Handle(request)
	return nil
}

func (client *RtspClient) Disconnect() {
	teardownCmd := commands.RtspTeardownCommand{
		SessionID: client.context.SessionHeader.Id,
	}
	response, sendErr := client.SendCommand(teardownCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		log.Panicln("Options failed!")
		return
	}
	client.connection.Close()
	// return nil
}

func (client *RtspClient) SendCommand(rtspCommand commands.RtspCommand) (*responses.RtspResponse, error) {
	if client.connection == nil {
		log.Panicln("Not connected!")
		return nil, errors.New("not connected")
	}
	// address := client.getAddress()
	// if rtspCommand.GetCommandType() == commands.RtspSetup && client.controlUri != nil {
	// 	address = *client.controlUri
	// }
	client.context.CSeq = client.context.CSeq + 1
	commandBuilder := commands.RtspCommandBuilder{
		Cseq:        client.context.CSeq,
		Address:     fmt.Sprintf("rtsp://%s:%d/%s", client.params.IP, client.params.Port, client.params.Path),
		AuthHeader:  client.context.AuthHeader,
		RtspCommand: rtspCommand,
		Credentials: client.params.Credentials,
		// ContentBase: client.contentBase,
	}
	log.Println(commandBuilder.BuildString())
	_, err := fmt.Fprintf(client.connection, commandBuilder.BuildString())
	if err != nil {
		return nil, err
	}
	response := <-client.rtspConnReader.RtspResponseChan
	log.Println(response.OriginalString)
	if response.StatusCode == 401 && len(response.AuthHeaders) > 0 && client.context.AuthHeader == nil {
		client.context.AuthHeader = response.AuthHeaders[0]
		return client.SendCommand(rtspCommand)
	}
	return &response, nil
}

func (client *RtspClient) GetContext() *client.RtspConnectionContext {
	return &client.context
}

func (client *RtspClient) exposeData(packetChan chan data.DataPacket) {
	for p := range packetChan {
		client.DataPacketChan <- p
	}
}

func (client *RtspClient) StartListening(transportHeader headers.TransportHeader) {
	if transportHeader.TransportType == headers.RtpAvpTcp {
		go client.exposeData(client.rtspConnReader.DataPacketChan)
	}
}
