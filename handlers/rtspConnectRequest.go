package handlers

import (
	"github.com/ardroh/goRtspClient"
	"github.com/ardroh/goRtspClient/commands"
)

type RtspConnectRequest struct {
	RtspClient       goRtspClient.RtspClient
	AvailableMethods []commands.RtspCommandTypes
}

func (request RtspConnectRequest) HasMethod(methodNameToFind commands.RtspCommandTypes) bool {
	for _, v := range request.AvailableMethods {
		if v == methodNameToFind {
			return true
		}
	}
	return false
}
