package responses

import (
	"sync"
)

type RtspResponseBuffer struct {
	responses []RtspResponse
	mutex     sync.Mutex
}

func (buffer *RtspResponseBuffer) Append(rtspResponse RtspResponse) {
	buffer.mutex.Lock()
	buffer.responses = append(buffer.responses, rtspResponse)
	buffer.mutex.Unlock()
}

func (buffer *RtspResponseBuffer) Get(cseq int) *RtspResponse {
	buffer.mutex.Lock()
	var foundResponse *RtspResponse
	var idxOfFound = -1
	for idx, e := range buffer.responses {
		if e.Cseq == cseq {
			foundResponse = &e
			idxOfFound = idx
			break
		}
	}
	if idxOfFound != -1 {
		buffer.responses = append(buffer.responses[:idxOfFound], buffer.responses[idxOfFound+1:]...)
	}
	buffer.mutex.Unlock()
	return foundResponse
}
