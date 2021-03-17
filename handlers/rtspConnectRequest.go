package handlers

import (
	"github.com/ardroh/goRtspClient/client"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/sdp"
)

type RtspConnectRequest struct {
	RtspClient       client.RtspClientIfc
	AvailableMethods []commands.RtspCommandTypes
	Sdp              sdp.Sdp
	Transmission     headers.RtspTransmissionType
	Transport        headers.RtspTransportType
	Session          headers.SessionHeader
}

func (request RtspConnectRequest) HasMethod(methodNameToFind commands.RtspCommandTypes) bool {
	for _, v := range request.AvailableMethods {
		if v == methodNameToFind {
			return true
		}
	}
	return false
}
