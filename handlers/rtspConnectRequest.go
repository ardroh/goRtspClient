package handlers

import (
	"github.com/ardroh/goRtspClient"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/sdp"
)

type RtspConnectRequest struct {
	RtspClient       goRtspClient.RtspClient
	AvailableMethods []commands.RtspCommandTypes
	Sdp              sdp.Sdp
	Transmission     headers.RtspTransmissionType
	Transport        headers.RtspTransportType
}

func (request RtspConnectRequest) HasMethod(methodNameToFind commands.RtspCommandTypes) bool {
	for _, v := range request.AvailableMethods {
		if v == methodNameToFind {
			return true
		}
	}
	return false
}
