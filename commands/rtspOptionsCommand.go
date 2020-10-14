package commands

type RtspOptionsCommand struct {
}

func (cmd RtspOptionsCommand) GetCommandType() RtspCommandType {
	return RtspOptions
}

func (cmd RtspOptionsCommand) GetParamsString() string {
	return "User-Agent: LibVLC/3.0.5 (LIVE555 Streaming Media v2016.11.28)\n"
}
