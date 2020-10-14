package commands

type RtspDescribeCommand struct {
}

func (cmd RtspDescribeCommand) GetCommandType() RtspCommandType {
	return RtspDescribe
}

func (cmd RtspDescribeCommand) GetParamsString() string {
	return "Accept: application/sdp\n"
}
