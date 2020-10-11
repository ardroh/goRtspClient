package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type RtpPacket struct {
	buffer []byte
	size   int
}

type RtspClient struct {
	cSeq          int
	rtspPath      string
	ip            string
	port          int
	connection    net.Conn
	readPacket    chan RtpPacket
	sessionID     string
	timeout       int
	lastKeepAlive time.Time
}

func (client *RtspClient) connect() bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ip, client.port))
	if err != nil {
		log.Panicln(err)
		return false
	}
	client.connection = conn
	optionsCmd := RtspOptionsCommand{
		cseq:    client.getNextCSeq(),
		address: client.getAddress(),
	}
	response, sendErr := client.send(optionsCmd)
	if sendErr != nil || response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return false
	}
	optionsResp := RtspOptionsResponse{
		rtspResponse: *response,
	}
	if !optionsResp.isMethodAvailable(Describe) {
		return false
	}
	describeCmd := RtspDescribeCommand{
		address: client.getAddress(),
		cseq:    client.getNextCSeq(),
	}
	response, sendErr = client.send(describeCmd)
	if sendErr != nil || response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return false
	}
	setupCmd := RtspSetupCommand{
		address:      client.getAddress(),
		cseq:         client.getNextCSeq(),
		transport:    RtpAvpTcp,
		transmission: Unicast,
		interleavedPair: InterleavedPair{
			rangeMin: 0,
			rangeMax: 1,
		},
	}
	response, sendErr = client.send(setupCmd)
	if sendErr != nil || response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return false
	}
	setupResp := RtspSetupResponse{
		rtspResponse: *response,
	}
	client.sessionID = setupResp.getSession()
	client.timeout = setupResp.getTimeout()
	playCmd := RtspPlayCommand{
		address:   client.getAddress(),
		cseq:      client.getNextCSeq(),
		sessionID: setupResp.getSession(),
	}
	response, sendErr = client.send(playCmd)
	if sendErr != nil || response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return false
	}
	client.readPacket = make(chan RtpPacket)
	go client.startReading()
	go client.keepAlive()
	return true
}

func (client *RtspClient) disconnect() error {
	teardownCmd := RtspTeardownCommand{
		address:   client.getAddress(),
		cseq:      client.getNextCSeq(),
		sessionID: client.sessionID,
	}
	response, sendErr := client.send(teardownCmd)
	if sendErr != nil || response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return sendErr
	}
	client.connection.Close()
	return nil
}

func (client *RtspClient) getNextCSeq() int {
	client.cSeq++
	return client.cSeq
}

func (client *RtspClient) getAddress() string {
	return fmt.Sprintf("rtsp://%s:%d/%s", client.ip, client.port, client.rtspPath)
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

func (client *RtspClient) send(rtspCommand RtspCommand) (*RtspResponse, error) {
	if client.connection == nil {
		log.Panicln("Not connected!")
		return nil, errors.New("not connected")
	}
	log.Println(rtspCommand.String())
	_, err := fmt.Fprintf(client.connection, rtspCommand.String())
	if err != nil {
		return nil, err
	}
	responseChan := make(chan string)
	go readResponse(client.connection, responseChan)
	response := RtspResponse{
		OriginalString: <-responseChan,
	}
	log.Println(response.OriginalString)
	return &response, nil
}

func (client *RtspClient) startReading() {
	reader := bufio.NewReader(client.connection)
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			return
		}
		client.readPacket <- RtpPacket{
			buffer: buffer,
			size:   bytesRead,
		}
	}
}

func (client *RtspClient) keepAlive() {
	client.lastKeepAlive = time.Now()
	for {
		t := time.Now()
		elapsed := t.Sub(client.lastKeepAlive)
		if elapsed < time.Duration(client.timeout/2)*time.Second {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		optionsCmd := RtspOptionsCommand{
			cseq:    client.getNextCSeq(),
			address: client.getAddress(),
		}
		_, err := client.send(optionsCmd)
		if err != nil {
			log.Println("Error on send keepalive. Exiting loop.")
			return
		}
		client.lastKeepAlive = time.Now()
	}
}
