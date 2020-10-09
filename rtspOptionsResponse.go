package main

import (
	"log"
	"regexp"
	"strings"
)

type RtspOptionsResponse struct {
	rtspResponse RtspResponse
}

func (options *RtspOptionsResponse) getAvailableCmds() []RtspCommandTypes {
	r, _ := regexp.Compile("Public: (.*?)\r\n")
	matches := r.FindStringSubmatch(options.rtspResponse.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get status code!")
		return []RtspCommandTypes{}
	}
	literals := strings.Split(matches[1], ",")
	var availableCmds []RtspCommandTypes
	for _, literal := range literals {
		availableCmds = append(availableCmds, RtspCommandTypes(strings.TrimSpace(literal)))
	}
	return availableCmds
}

func (options *RtspOptionsResponse) isMethodAvailable(cmdType RtspCommandTypes) bool {
	for _, t := range options.getAvailableCmds() {
		if t == cmdType {
			return true
		}
	}
	return false
}
