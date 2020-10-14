package commands

type RtspOptionsCommand struct {
}

func (cmd RtspOptionsCommand) GetCommandType() RtspCommandType {
	return RtspOptions
}

func (cmd RtspOptionsCommand) GetParamsString() string {
	return ""
}
