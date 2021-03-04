package parsers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ardroh/goRtspClient/responses"
)

func getNumberFromLine(pattern string, line string) (int, error) {
	r, _ := regexp.Compile(pattern)
	matches := r.FindStringSubmatch(line)
	num, err := strconv.Atoi(strings.TrimSpace(matches[1]))
	return num, err
}

type RtspResponseParser struct {
}

func (parser *RtspResponseParser) Parse(command string) (*responses.RtspResponse, error) {
	var cseq int
	var statusCode responses.RtspResponseCodes
	lines := strings.Split(command, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "CSeq:") {
			num, err := getNumberFromLine("CSeq: (.*)", line)
			if err != nil {
				return nil, err
			}
			cseq = num
		} else if strings.HasPrefix(line, "RTSP/1.0") {
			num, err := getNumberFromLine("RTSP/1.0 (.*?) (.*)", line)
			if err != nil {
				return nil, err
			}
			statusCode = responses.RtspResponseCodes(num)
		}
	}
	return responses.CreateRtspResponse(command, cseq, statusCode), nil
}
