package parsers

import (
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

type RtspResponseParser struct {
}

func (parser RtspResponseParser) Parse(responseLiteral string) (responses.RtspResponse, error) {
	response := responses.RtspResponse{
		OriginalString: responseLiteral,
	}
	lines := strings.Split(responseLiteral, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RTSP/1.0") {
			num := getNumberFromLine(line, "RTSP/1.0 (.*?) (.*)")
			response.StatusCode = responses.RtspResponseCodes(num)
		} else if strings.HasPrefix(line, "CSeq:") {
			response.Cseq = getNumberFromLine(line, "CSeq: (.*)")
		} else if strings.HasPrefix(line, "Content-Type") {
			response.ContentType = getStringFromLine(line, "Content-Type: (.*)")
		} else if strings.HasPrefix(line, "Content-Base") {
			response.ContentBase = getStringFromLine(line, "Content-Base: (.*)")
		} else if strings.HasPrefix(line, "Content-Length:") {
			response.ContentLength = getNumberFromLine(line, "Content-Length: (.*)")
		} else if strings.HasPrefix(line, "Server:") {
			response.Server = getStringFromLine(line, "Server: (.*)")
		} else if strings.HasPrefix(line, "Date:") {
			response.DateTime = getRtspDateParserFromLine(line, "Date: (.*)")
		}
	}
	return response, nil
}
