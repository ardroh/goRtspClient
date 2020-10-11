package main

import (
	"fmt"
	"strings"
)

type RtspTeardownCommand struct {
	address   string
	cseq      int
	sessionID string
}

func (teardown RtspTeardownCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TEARDOWN %s RTSP/1.0\n", teardown.address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", teardown.cseq))
	sb.WriteString(fmt.Sprintf("Session: %s\n\n", teardown.sessionID))
	return sb.String()
}
