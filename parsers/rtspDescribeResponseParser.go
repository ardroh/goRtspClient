package parsers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

type RtspDescribeResponseParser struct {
}

func (parser RtspDescribeResponseParser) FromString(literal string) (*responses.RtspDescribeResponse, error) {
	baseResponse, err := RtspResponseParser{}.Parse(literal)
	if err != nil {
		return nil, errors.New("Failed to parse base response")
	}
	return parser.FromBaseResponse(baseResponse)
}

func (parser RtspDescribeResponseParser) FromBaseResponse(rtspResponse responses.RtspResponse) (*responses.RtspDescribeResponse, error) {
	response := &responses.RtspDescribeResponse{RtspResponse: rtspResponse}
	r, err := regexp.Compile("[\t ]*([a-z]=.*)")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(response.RtspResponse.OriginalString, "\n")
	var mediaParams []string
	for _, l := range lines {
		param := r.FindStringSubmatch(l)
		if len(param) == 2 {
			mediaParams = append(mediaParams, param[1])
		}
	}
	if len(mediaParams) < 2 && rtspResponse.ContentLength > 0 {
		return nil, errors.New("Failed to parse content")
	}
	content := strings.Join(mediaParams[:], "\n")
	content = strings.TrimSpace(content)
	mediaInfo, err := SdpParser{}.FromString(content)
	if mediaInfo == nil || err != nil {
		return nil, err
	}
	response.Sdp = *mediaInfo
	return response, nil
}
