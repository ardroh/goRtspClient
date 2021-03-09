package responses

import (
	"regexp"
	"strings"

	"github.com/ardroh/goRtspClient/sdp"
)

type RtspDescribeResponse struct {
	RtspResponse RtspResponse
	Sdp          sdp.Sdp
}

func InitRtspDescribeResponse(rtspResponse RtspResponse) *RtspDescribeResponse {
	return &RtspDescribeResponse{
		RtspResponse: rtspResponse,
	}
}

func (resp RtspDescribeResponse) GetControlUris() []string {
	var uris []string
	lines := strings.Split(resp.RtspResponse.OriginalString, "\n")
	for _, line := range lines {
		r, _ := regexp.Compile("a=control:(.*)")
		matches := r.FindStringSubmatch(line)
		if len(matches) < 1 {
			continue
		}
		address := strings.TrimSpace(matches[1])
		uris = append(uris, address)
	}
	return uris
}
