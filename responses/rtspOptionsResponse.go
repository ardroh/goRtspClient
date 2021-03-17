package responses

import (
	"github.com/ardroh/goRtspClient/commands"
)

// RtspOptionsResponse response parser
type RtspOptionsResponse struct {
	RtspResponse     RtspResponse
	AvailableMethods []commands.RtspCommandTypes
}

func (options *RtspOptionsResponse) IsMethodAvailable(cmdType commands.RtspCommandTypes) bool {
	for _, t := range options.AvailableMethods {
		if t == cmdType {
			return true
		}
	}
	return false
}
