package parsers

import (
	"errors"
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
	splitHeaderAndContent := strings.Split(rtspResponse.OriginalString, "\n\n")
	if len(splitHeaderAndContent) < 2 && rtspResponse.ContentLength != 0 {
		return nil, errors.New("Failed to parse content")
	}
	content := splitHeaderAndContent[1]
	content = strings.TrimSpace(content)
	mediaInfo, err := MediaInformationParser{}.FromString(content)
	if mediaInfo == nil || err != nil {
		return nil, err
	}
	response.MediaInfo = *mediaInfo
	return response, nil
}
