package commands

import (
	"fmt"
	"strings"

	"github.com/ardroh/goRtspClient/auth"
)

type RtspCommandBuilder struct {
	Cseq        int
	Address     string
	AuthHeader  auth.RtspAuthHeader
	RtspCommand RtspCommand
}

func (builder RtspCommandBuilder) BuildString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s RTSP/1.0\n", builder.RtspCommand.GetCommandType(), builder.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", builder.Cseq))
	if builder.AuthHeader != nil {
		sb.WriteString(builder.AuthHeader.String())
	}
	sb.WriteString(builder.RtspCommand.GetParamsString())
	sb.WriteString("\n\n")
	return sb.String()
}
