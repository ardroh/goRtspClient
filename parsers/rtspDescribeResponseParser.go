package parsers

import (
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

type RtspDescribeResponseParser struct {
}

func (parser *RtspDescribeResponseParser) Parse(rtspResponse *responses.RtspResponse) *responses.RtspDescribeResponse {
	response := &responses.RtspDescribeResponse{RtspResponse: *rtspResponse}
	lines := strings.Split(rtspResponse.OriginalString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
	}
	return response
}
