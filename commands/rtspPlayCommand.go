package commands

import (
	"fmt"
)

type RtspPlayCommand struct {
	SessionID string
}

func (cmd RtspPlayCommand) GetCommandType() RtspCommandType {
	return RtspPlay
}

func (cmd RtspPlayCommand) GetParamsString() string {
	return fmt.Sprintf("Session: %s\n\n", cmd.SessionID)
}
