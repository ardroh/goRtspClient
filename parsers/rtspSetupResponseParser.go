package parsers

import (
	"errors"
	"strings"

	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/responses"
)

type RtspSetupResponseParser struct {
}

func (parser RtspSetupResponseParser) FromString(literal string) (responses.RtspSetupResponse, error) {
	baseResponse, err := RtspResponseParser{}.Parse(literal)
	if err != nil {
		return responses.RtspSetupResponse{}, errors.New("Failed to parse base response")
	}
	return parser.FromBaseResponse(baseResponse)
}

func (parser RtspSetupResponseParser) FromBaseResponse(rtspResponse responses.RtspResponse) (responses.RtspSetupResponse, error) {
	response := responses.RtspSetupResponse{RtspResponse: rtspResponse}
	lines := strings.Split(rtspResponse.OriginalString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		var err error
		if strings.HasPrefix(line, "Session:") {
			response.SessionHeader, err = headers.SessionHeader{}.FromString(line)
			if err != nil {
				return response, errors.New("Can't parse session info!")
			}
		} else if strings.HasPrefix(line, "Transport:") {
			response.Transport, err = headers.TransportHeader{}.FromString(line)
			if err != nil {
				return response, errors.New("Can't parse transport header!")
			}
		}
	}
	return response, nil
}
