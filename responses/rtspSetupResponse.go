package responses

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type RtspSetupResponse struct {
	rtspResponse RtspResponse
}

func InitRtspSetupResponse(rtspResponse RtspResponse) *RtspSetupResponse {
	return &RtspSetupResponse{
		rtspResponse: rtspResponse,
	}
}

func (setupResp RtspSetupResponse) GetSsrc() string {
	r, _ := regexp.Compile("Transport:(.+?);ssrc=(.*?);")
	matches := r.FindStringSubmatch(setupResp.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get ssrc!")
		return ""
	}
	return matches[2]
}

func (setupResp RtspSetupResponse) GetTimeout() int {
	r, _ := regexp.Compile("Session:(.+?)timeout=(.*)")
	matches := r.FindStringSubmatch(setupResp.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get session!")
		return -1
	}
	num, err := strconv.Atoi(strings.TrimSpace(matches[2]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return num
}

func (setupResp RtspSetupResponse) GetSession() string {
	r, _ := regexp.Compile("Session: (.+?);")
	matches := r.FindStringSubmatch(setupResp.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get session!")
		return ""
	}
	return matches[1]
}
