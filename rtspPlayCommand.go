package goRtspClient

import (
	"fmt"
	"strings"
)

type RtspPlayCommand struct {
	address   string
	cseq      int
	sessionID string
}

func (play RtspPlayCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PLAY %s RTSP/1.0\n", play.address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", play.cseq))
	sb.WriteString(fmt.Sprintf("Session: %s\n\n", play.sessionID))
	return sb.String()
}
