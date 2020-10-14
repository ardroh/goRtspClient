package goRtspClient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ardroh/goRtspClient/auth"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/responses"
	"github.com/ardroh/goRtspClient/rtp"
)

type RtspConnectionParams struct {
	IP           string
	Port         int
	Path         string
	Transmission commands.RtspTransmissionType
	Transport    commands.RtspTransportType
	Credentials  auth.Credentials
}

type rtspClient struct {
	cSeq             int
	connection       net.Conn
	readPacket       chan rtp.RtpPacket
	sessionID        string
	timeout          int
	lastKeepAlive    time.Time
	connectionParams RtspConnectionParams
	authHeader       auth.RtspAuthHeader
}

func InitRtspClient(params RtspConnectionParams) *rtspClient {
	return &rtspClient{
		connectionParams: params,
	}
}

func (client *rtspClient) GetReadChan() chan rtp.RtpPacket {
	return client.readPacket
}

func (client *rtspClient) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.connectionParams.IP, client.connectionParams.Port))
	if err != nil {
		return err
	}
	client.connection = conn
	optionsCmd := commands.RtspOptionsCommand{}
	response, sendErr := client.send(optionsCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		return sendErr
	}
	optionsResp := responses.InitRtspOptionsResponse(*response)
	if !optionsResp.IsMethodAvailable(commands.Describe) {
		return sendErr
	}
	describeCmd := commands.RtspDescribeCommand{}
	response, sendErr = client.send(describeCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		return sendErr
	}
	if client.connectionParams.Transport != commands.RtpAvpTcp || client.connectionParams.Transmission != commands.Unicast {
		return errors.New("unsupported transport or transmission")
	}
	setupCmd := commands.RtspSetupCommand{
		Transport:    client.connectionParams.Transport,
		Transmission: client.connectionParams.Transmission,
		InterleavedPair: commands.InterleavedPair{
			RangeMin: 0,
			RangeMax: 1,
		},
	}
	response, sendErr = client.send(setupCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		return sendErr
	}
	setupResp := responses.InitRtspSetupResponse(*response)
	client.sessionID = setupResp.GetSession()
	client.timeout = setupResp.GetTimeout()
	playCmd := commands.RtspPlayCommand{
		SessionID: setupResp.GetSession(),
	}
	response, sendErr = client.send(playCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		return sendErr
	}
	client.readPacket = make(chan rtp.RtpPacket)
	go client.startReading()
	go client.keepAlive()
	return nil
}

func (client *rtspClient) Disconnect() error {
	teardownCmd := commands.RtspTeardownCommand{
		SessionID: client.sessionID,
	}
	response, sendErr := client.send(teardownCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		log.Panicln("Options failed!")
		return sendErr
	}
	client.connection.Close()
	return nil
}

func (client *rtspClient) getNextCSeq() int {
	client.cSeq++
	return client.cSeq
}

func (client *rtspClient) getAddress() string {
	return fmt.Sprintf("rtsp://%s:%d/%s", client.connectionParams.IP, client.connectionParams.Port, client.connectionParams.Path)
}

func readResponse(conn net.Conn, responseChan chan string) {
	i := 0
	reader := bufio.NewReader(conn)
	var sb strings.Builder
	for i < 100 {
		i++
		lineRead, err := reader.ReadString('\n')
		if err == nil {
			sb.WriteString(lineRead)
		}
		lineLength := len(lineRead)
		if lineLength == 2 {
			responseChan <- sb.String()
			break
		}
	}
	close(responseChan)
}

func (client *rtspClient) send(rtspCommand commands.RtspCommand) (*responses.RtspResponse, error) {
	if client.connection == nil {
		log.Panicln("Not connected!")
		return nil, errors.New("not connected")
	}
	commandBuilder := commands.RtspCommandBuilder{
		Cseq:        client.getNextCSeq(),
		Address:     client.getAddress(),
		AuthHeader:  client.authHeader,
		RtspCommand: rtspCommand,
	}
	log.Println(commandBuilder.BuildString())
	_, err := fmt.Fprintf(client.connection, commandBuilder.BuildString())
	if err != nil {
		return nil, err
	}
	responseChan := make(chan string)
	go readResponse(client.connection, responseChan)
	response := responses.RtspResponse{
		OriginalString: <-responseChan,
	}
	if response.GetStatusCode() == responses.RtspUnauthorized {
		authRequest := response.GetRtspAuthType()
		client.authHeader = auth.BuildRtspAuthHeader(authRequest, client.connectionParams.Credentials)
		return client.send(rtspCommand) //retry
	}
	log.Println(response.OriginalString)
	return &response, nil
}

func (client *rtspClient) startReading() {
	reader := bufio.NewReader(client.connection)
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			return
		}
		client.readPacket <- rtp.RtpPacket{
			Buffer: buffer,
			Size:   bytesRead,
		}
	}
}

func (client *rtspClient) keepAlive() {
	client.lastKeepAlive = time.Now()
	for {
		t := time.Now()
		elapsed := t.Sub(client.lastKeepAlive)
		if elapsed < time.Duration(client.timeout/2)*time.Second {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		optionsCmd := commands.RtspOptionsCommand{}
		_, err := client.send(optionsCmd)
		if err != nil {
			log.Println("Error on send keepalive. Exiting loop.")
			return
		}
		client.lastKeepAlive = time.Now()
	}
}
