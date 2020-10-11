package main

import (
	"log"
	"regexp"
)

type RtspSetupResponse struct {
	rtspResponse RtspResponse
}

func (setupResp RtspSetupResponse) getSsrc() string {
	r, _ := regexp.Compile("Transport:(.+?);ssrc=(.*?);")
	matches := r.FindStringSubmatch(setupResp.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get ssrc!")
		return ""
	}
	return matches[2]
}

func (setupResp RtspSetupResponse) getSession() string {
	r, _ := regexp.Compile("Session: (.+?);")
	matches := r.FindStringSubmatch(setupResp.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get session!")
		return ""
	}
	return matches[1]
}
