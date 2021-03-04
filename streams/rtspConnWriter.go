package streams

import (
	"fmt"
	"net"
	"sync"
)

type rtspConnWriter struct {
	conn  net.Conn
	mutex *sync.Mutex
}

func CreateRtspConnWriter(conn net.Conn, mutex *sync.Mutex) *rtspConnWriter {
	return &rtspConnWriter{
		conn:  conn,
		mutex: mutex,
	}
}

func (writer *rtspConnWriter) Send(message string) error {
	writer.mutex.Lock()
	_, err := fmt.Fprintf(writer.conn, message)
	writer.mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}
