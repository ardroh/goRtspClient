package parsers

import (
	"errors"
	"strings"

	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/responses"
)

type RtspOptionsResponseParser struct {
}

func (parser RtspOptionsResponseParser) FromString(literal string) (*responses.RtspOptionsResponse, error) {
	baseResponse, err := RtspResponseParser{}.Parse(literal)
	if err != nil {
		return nil, errors.New("Failed to parse base response")
	}
	return parser.FromBaseResponse(baseResponse)
}

func (parser RtspOptionsResponseParser) FromBaseResponse(baseResposne responses.RtspResponse) (*responses.RtspOptionsResponse, error) {
	response := &responses.RtspOptionsResponse{
		RtspResponse: baseResposne,
	}
	lines := strings.Split(baseResposne.OriginalString, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Public:") {
			methods := getStringFromLine(line, "Public: (.*)")
			if methods == nil {
				continue
			}
			literals := strings.Split(*methods, ",")
			for _, literal := range literals {
				response.AvailableMethods = append(response.AvailableMethods, commands.RtspCommandTypes(strings.TrimSpace(literal)))
			}
		}
	}
	return response, nil
}
