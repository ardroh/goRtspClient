package commands

import (
	"fmt"
)

type RtspTeardownCommand struct {
	SessionID string
}

func (cmd RtspTeardownCommand) GetCommandType() RtspCommandType {
	return RtspTeardown
}

func (cmd RtspTeardownCommand) GetParamsString() string {
	return fmt.Sprintf("Session: %s\n\n", cmd.SessionID)
}
