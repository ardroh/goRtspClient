package goRtspClient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/ardroh/goRtspClient/auth"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/parsers"
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

type RtspClient interface {
	SendCommand(rtspCommand commands.RtspCommand) (*responses.RtspResponse, error)
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
	contentBase      *string
	controlUri       *string
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
	optionsResp, err := parsers.RtspOptionsResponseParser{}.FromBaseResponse(*response)
	if !optionsResp.IsMethodAvailable(commands.Describe) {
		return sendErr
	}
	describeCmd := commands.RtspDescribeCommand{}
	response, sendErr = client.send(describeCmd)
	if sendErr != nil || response.GetStatusCode() != responses.RtspOk {
		return sendErr
	}
	describeResp := responses.InitRtspDescribeResponse(*response)
	if describeResp.RtspResponse.ContentBase != nil {
		client.contentBase = describeResp.RtspResponse.ContentBase
	}
	controlUris := describeResp.GetControlUris()
	if len(controlUris) > 0 {
		client.controlUri = &controlUris[len(controlUris)-1]
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
	setupResp, parseError := parsers.RtspSetupResponseParser{}.FromBaseResponse(*response)
	if parseError != nil {
		return parseError
	}
	client.sessionID = setupResp.SessionInfo.Id
	client.timeout = setupResp.SessionInfo.Timeout
	playCmd := commands.RtspPlayCommand{
		SessionID: setupResp.SessionInfo.Id,
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

func peekIsRtspMessage(reader *bufio.Reader) (bool, error) {
	peekedBytes, err := reader.Peek(8)
	if err != nil {
		return false, err
	}
	peekedLine := string(peekedBytes[:])
	return peekedLine == "RTSP/1.0", nil
}

func readResponse(conn net.Conn) (*string, error) {
	reader := bufio.NewReader(conn)
	for {
		isRtspMessage, err := peekIsRtspMessage(reader)
		if err != nil {
			return nil, err
		}
		if !isRtspMessage {
			continue
		}
		buffer := make([]byte, 2048)
		len, err := reader.Read(buffer)
		if err != nil {
			return nil, err
		}
		literalData := string(buffer[:len])
		return &literalData, nil
	}
}

func (client *rtspClient) send(rtspCommand commands.RtspCommand) (*responses.RtspResponse, error) {
	if client.connection == nil {
		log.Panicln("Not connected!")
		return nil, errors.New("not connected")
	}
	address := client.getAddress()
	if rtspCommand.GetCommandType() == commands.RtspSetup && client.controlUri != nil {
		address = *client.controlUri
	}
	commandBuilder := commands.RtspCommandBuilder{
		Cseq:        client.getNextCSeq(),
		Address:     address,
		AuthHeader:  client.authHeader,
		RtspCommand: rtspCommand,
		ContentBase: client.contentBase,
	}
	log.Println(commandBuilder.BuildString())
	_, err := fmt.Fprintf(client.connection, commandBuilder.BuildString())
	if err != nil {
		return nil, err
	}
	responseString, err := readResponse(client.connection)
	if err != nil {
		return nil, err
	}
	response := responses.RtspResponse{
		OriginalString: *responseString,
	}
	log.Println(response.OriginalString)
	if response.GetStatusCode() == responses.RtspUnauthorized {
		authRequest := response.GetRtspAuthType()
		client.authHeader = auth.BuildRtspAuthHeader(authRequest, client.connectionParams.Credentials)
		return client.send(rtspCommand) //retry
	}
	return &response, nil
}

func (client *rtspClient) startReading() {
	reader := bufio.NewReader(client.connection)
	for {
		isRtspMessage, err := peekIsRtspMessage(reader)
		if err != nil {
			return
		}
		if isRtspMessage {
			continue
		}
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
