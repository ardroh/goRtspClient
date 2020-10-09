package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type RtspClient struct {
	cSeq       int
	rtspPath   string
	ip         string
	port       int
	connection net.Conn
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
	response := client.send(optionsCmd)
	if response == nil || response.getStatusCode() != RtspOk {
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
		cseq:    client.getNextCSeq(),
		address: client.getAddress(),
	}
	response = client.send(describeCmd)
	if response == nil || response.getStatusCode() != RtspOk {
		log.Panicln("Options failed!")
		return false
	}
	return true
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

func (client *RtspClient) send(rtspCommand RtspCommand) *RtspResponse {
	if client.connection == nil {
		log.Panicln("Not connected!")
		return nil
	}
	log.Println(rtspCommand.String())
	fmt.Fprintf(client.connection, rtspCommand.String())
	responseChan := make(chan string)
	go readResponse(client.connection, responseChan)
	response := RtspResponse{
		OriginalString: <-responseChan,
	}
	log.Println(response.OriginalString)
	return &response
}
