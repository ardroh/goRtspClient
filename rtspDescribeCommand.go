package goRtspClient

import (
	"fmt"
	"strings"
)

type RtspDescribeCommand struct {
	address string
	cseq    int
}

func (cmd RtspDescribeCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DESCRIBE %s RTSP/1.0\n", cmd.address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", cmd.cseq))
	sb.WriteString("Accept: application/sdp\n\n")
	return sb.String()
}
