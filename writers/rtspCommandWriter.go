package writers

import (
	"fmt"
	"net"
	"sync"
)

type rtspCommandWriter struct {
	conn  net.Conn
	mutex *sync.Mutex
}

func CreateRtspConnWriter(conn net.Conn, mutex *sync.Mutex) *rtspCommandWriter {
	return &rtspCommandWriter{
		conn:  conn,
		mutex: mutex,
	}
}

func (writer *rtspCommandWriter) Send(message string) error {
	writer.mutex.Lock()
	_, err := fmt.Fprintf(writer.conn, message)
	writer.mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}
