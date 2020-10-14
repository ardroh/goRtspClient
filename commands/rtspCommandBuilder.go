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
	ContentBase *string
}

func (builder RtspCommandBuilder) BuildString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s RTSP/1.0\n", builder.RtspCommand.GetCommandType(), builder.Address))
	sb.WriteString(fmt.Sprintf("CSeq: %d\n", builder.Cseq))
	sb.WriteString("User-Agent: goRtspClient\n")
	if builder.AuthHeader != nil {
		authAddress := builder.Address
		if builder.ContentBase != nil {
			authAddress = *builder.ContentBase
		}
		sb.WriteString(builder.AuthHeader.String(string(builder.RtspCommand.GetCommandType()), authAddress) + "\n")
	}
	sb.WriteString(builder.RtspCommand.GetParamsString())
	sb.WriteString("\n")
	return sb.String()
}
