package parsers

import (
	"errors"
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

type RtspSetupResponseParser struct {
}

func (parser RtspSetupResponseParser) FromString(literal string) (*responses.RtspSetupResponse, error) {
	baseResponse, err := RtspResponseParser{}.Parse(literal)
	if err != nil {
		return nil, errors.New("Failed to parse base response")
	}
	return parser.FromBaseResponse(baseResponse)
}

func (parser RtspSetupResponseParser) FromBaseResponse(rtspResponse responses.RtspResponse) (*responses.RtspSetupResponse, error) {
	response := &responses.RtspSetupResponse{RtspResponse: rtspResponse}
	lines := strings.Split(rtspResponse.OriginalString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Session:") {
			sessionInfo := parseSessionInfo(line)
			if sessionInfo == nil {
				return nil, errors.New("Can't parse session info!")
			}
			response.SessionInfo = *sessionInfo
		}
	}
	return response, nil
}

func parseSessionInfo(line string) (sessionInfo *responses.SessionInfo) {
	sessionTimeout := getNumberFromLine(line, "Session: .*? timeout=(.*)")
	sessionId := getStringFromLine(line, "Session: (.*);.*")
	if sessionId == nil {
		return nil
	}
	return &responses.SessionInfo{
		Id:      *sessionId,
		Timeout: sessionTimeout,
	}
}
