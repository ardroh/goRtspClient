package readers

import (
	"bufio"
	"log"
	"net"

	"github.com/ardroh/goRtspClient/parsers"
	"github.com/ardroh/goRtspClient/responses"
	"github.com/ardroh/goRtspClient/rtp"
)

func peekIsRtspMessage(reader *bufio.Reader) (bool, error) {
	peekedBytes, err := reader.Peek(8)
	if err != nil {
		return false, err
	}
	peekedLine := string(peekedBytes[:])
	return peekedLine == "RTSP/1.0", nil
}

type RtspSocketReader struct {
	conn             net.Conn
	RtpPacketChan    chan rtp.RtpPacket
	RtspResponseChan chan responses.RtspResponse
}

func CreateRtspConnReader(conn net.Conn) *RtspSocketReader {
	return &RtspSocketReader{
		conn:             conn,
		RtpPacketChan:    make(chan rtp.RtpPacket),
		RtspResponseChan: make(chan responses.RtspResponse),
	}
}

func getRtpDataStartIndex(buffer []byte) (bool, int) {
	for i := range buffer {
		if i+1 >= len(buffer) {
			continue
		}

		if buffer[i] == byte('$') &&
			buffer[i+1] == byte(0) {
			return true, i
		}
	}
	return false, 0
}

func (reader *RtspSocketReader) StartReading() {
	bufferReader := bufio.NewReader(reader.conn)
	buffer := make([]byte, 4096)
	for {
		for i := range buffer {
			buffer[i] = 0
		}
		isRtspMessage, err := peekIsRtspMessage(bufferReader)
		if err != nil {
			return
		}
		var len int
		len, err = bufferReader.Read(buffer)
		if len == 0 {
			continue
		}
		copiedBuffer := make([]byte, len)
		copy(copiedBuffer, buffer[:len])
		if isRtspMessage {
			hasRtpDataMixed, rtpStartIdx := getRtpDataStartIndex(copiedBuffer)
			if hasRtpDataMixed {
				reader.handleRtspData(copiedBuffer[:rtpStartIdx])
				reader.handleBinaryData(copiedBuffer[rtpStartIdx:])
			} else {
				reader.handleRtspData(copiedBuffer)
			}
		} else {
			reader.handleBinaryData(copiedBuffer)
		}
	}
}

func (reader *RtspSocketReader) handleBinaryData(buffer []byte) {
	reader.RtpPacketChan <- rtp.RtpPacket{
		Buffer: buffer,
		Size:   len(buffer),
	}
}

func (reader *RtspSocketReader) handleRtspData(buffer []byte) {
	literalData := string(buffer)
	response, err := parsers.RtspResponseParser{}.Parse(literalData)
	if err != nil {
		log.Fatalf("Failed to parse: %s", literalData)
		return
	}
	reader.RtspResponseChan <- response
}
