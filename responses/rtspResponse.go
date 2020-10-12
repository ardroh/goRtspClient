package responses

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type RtspResponse struct {
	OriginalString string
}

func (response *RtspResponse) GetCSeq() int {
	if len(response.OriginalString) == 0 {
		return -1
	}
	r, _ := regexp.Compile("CSeq: (.*?)\n")
	matches := r.FindStringSubmatch(response.OriginalString)
	num, err := strconv.Atoi(strings.TrimSpace(matches[1]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return num
}

func (response *RtspResponse) GetStatusCode() RtspResponseCodes {
	r, _ := regexp.Compile("(RTSP/1.0) (.*?) (.*)")
	matches := r.FindStringSubmatch(response.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get status code!")
		return -1
	}
	num, err := strconv.Atoi(strings.TrimSpace(matches[2]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return RtspResponseCodes(num)
}
