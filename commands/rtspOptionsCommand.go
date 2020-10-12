package commands

import (
	"fmt"
	"strings"
)

type RtspOptionsCommand struct {
	Address string
	Cseq    int
}

func (cmd RtspOptionsCommand) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OPTIONS %s RTSP/1.0\n", cmd.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", cmd.Cseq))
	sb.WriteString("User-Agent: LibVLC/3.0.5 (LIVE555 Streaming Media v2016.11.28)\n\n")
	return sb.String()
}
