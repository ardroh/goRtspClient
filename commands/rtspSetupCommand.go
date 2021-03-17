package commands

import (
	"strings"

	"github.com/ardroh/goRtspClient/headers"
)

type RtspSetupCommand struct {
	TransportHeader headers.TransportHeader
}

func (cmd RtspSetupCommand) GetCommandType() RtspCommandType {
	return RtspSetup
}

func (cmd RtspSetupCommand) GetParamsString() string {
	var sb strings.Builder
	sb.WriteString(cmd.TransportHeader.ToString())
	sb.WriteString("\n")
	return sb.String()
}
