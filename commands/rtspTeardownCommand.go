package commands

import (
	"fmt"
	"strings"
)

type RtspTeardownCommand struct {
	Address   string
	Cseq      int
	SessionID string
}

func (teardown RtspTeardownCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TEARDOWN %s RTSP/1.0\n", teardown.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", teardown.Cseq))
	sb.WriteString(fmt.Sprintf("Session: %s\n\n", teardown.SessionID))
	return sb.String()
}
