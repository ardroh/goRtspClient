package commands

import (
	"fmt"
	"strings"
)

type RtspDescribeCommand struct {
	Address string
	Cseq    int
}

func (cmd RtspDescribeCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DESCRIBE %s RTSP/1.0\n", cmd.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", cmd.Cseq))
	sb.WriteString("Accept: application/sdp\n\n")
	return sb.String()
}
