package parsers

import (
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

type RtspResponseParser struct {
}

func (parser *RtspResponseParser) Parse(responseLiteral string) (*responses.RtspResponse, error) {
	response := &responses.RtspResponse{
		OriginalString: responseLiteral,
	}
	lines := strings.Split(responseLiteral, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RTSP/1.0") {
			num := getNumberFromLine("RTSP/1.0 (.*?) (.*)", line)
			response.StatusCode = responses.RtspResponseCodes(num)
		} else if strings.HasPrefix(line, "CSeq:") {
			response.Cseq = getNumberFromLine("CSeq: (.*)", line)
		} else if strings.HasPrefix(line, "Content-Type") {
			response.ContentType = getStringFromLine("Content-Type: (.*)", line)
		} else if strings.HasPrefix(line, "Content-Base") {
			response.ContentBase = getStringFromLine("Content-Base: (.*)", line)
		} else if strings.HasPrefix(line, "Content-Length:") {
			response.ContentLength = getNumberFromLine("Content-Length: (.*)", line)
		} else if strings.HasPrefix(line, "Server:") {
			response.Server = getStringFromLine("Server: (.*)", line)
		} else if strings.HasPrefix(line, "Date:") {
			response.DateTime = getRtspDateParserFromLine("Date: (.*)", line)
		}
	}
	return response, nil
}
