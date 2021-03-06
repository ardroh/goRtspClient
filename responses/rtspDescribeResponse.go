package responses

import (
	"regexp"
	"strings"
)

type RtspDescribeResponse struct {
	RtspResponse RtspResponse
}

func InitRtspDescribeResponse(rtspResponse RtspResponse) *RtspDescribeResponse {
	return &RtspDescribeResponse{
		RtspResponse: rtspResponse,
	}
}

func (resp RtspDescribeResponse) GetContentBase() *string {
	r, _ := regexp.Compile("Content-Base: (.*)")
	matches := r.FindStringSubmatch(resp.RtspResponse.OriginalString)
	if len(matches) < 1 {
		return nil
	}
	contentBase := strings.TrimSpace(matches[1])
	return &contentBase
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
