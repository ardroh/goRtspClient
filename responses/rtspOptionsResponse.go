package responses

import (
	"log"
	"regexp"
	"strings"

	"github.com/ardroh/goRtspClient/commands"
)

// RtspOptionsResponse response parser
type RtspOptionsResponse struct {
	rtspResponse RtspResponse
}

// InitRtspOptionsResponse builds the RTSP/OPTIONS response parser
func InitRtspOptionsResponse(rtspResponse RtspResponse) *RtspOptionsResponse {
	return &RtspOptionsResponse{
		rtspResponse: rtspResponse,
	}
}

func (options *RtspOptionsResponse) GetAvailableCmds() []commands.RtspCommandTypes {
	r, _ := regexp.Compile("Public: (.*?)\r\n")
	matches := r.FindStringSubmatch(options.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get status code!")
		return []commands.RtspCommandTypes{}
	}
	literals := strings.Split(matches[1], ",")
	var availableCmds []commands.RtspCommandTypes
	for _, literal := range literals {
		availableCmds = append(availableCmds, commands.RtspCommandTypes(strings.TrimSpace(literal)))
	}
	return availableCmds
}

func (options *RtspOptionsResponse) IsMethodAvailable(cmdType commands.RtspCommandTypes) bool {
	for _, t := range options.GetAvailableCmds() {
		if t == cmdType {
			return true
		}
	}
	return false
}
