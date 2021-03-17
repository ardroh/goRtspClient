package readers

import (
	"bufio"
	"net"
	"sync"

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

type interleavedReader struct {
	conn             net.Conn
	RtpPacketChan    chan rtp.RtpPacket
	RtspResponseChan chan responses.RtspResponse
}

func CreateRtspConnReader(conn net.Conn, mutex *sync.Mutex) *interleavedReader {
	return &interleavedReader{
		conn: conn,
	}
}

func (reader *interleavedReader) StartReading() {
	bufferReader := bufio.NewReader(reader.conn)
	for {
		isRtspMessage, err := peekIsRtspMessage(bufferReader)
		if err != nil {
			return
		}
		var len int
		buffer := make([]byte, 2048)
		len, err = bufferReader.Read(buffer)
		if isRtspMessage {
			reader.handleRtspData(buffer, len)
		} else {
			reader.handleBinaryData(buffer, len)
		}
	}
}

func (reader *interleavedReader) handleBinaryData(buffer []byte, length int) {
	reader.RtpPacketChan <- rtp.RtpPacket{
		Buffer: buffer,
		Size:   length,
	}
}

func (reader *interleavedReader) handleRtspData(buffer []byte, length int) {
	literalData := string(buffer[:length])
	reader.RtspResponseChan <- responses.RtspResponse{
		OriginalString: literalData,
	}
}
