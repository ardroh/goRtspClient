package commands

import (
	"fmt"
	"strings"
)

type RtspPlayCommand struct {
	Address   string
	Cseq      int
	SessionID string
}

func (play RtspPlayCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PLAY %s RTSP/1.0\n", play.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", play.Cseq))
	sb.WriteString(fmt.Sprintf("Session: %s\n\n", play.SessionID))
	return sb.String()
}
