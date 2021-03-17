package client

import (
	"github.com/ardroh/goRtspClient/auth"
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/responses"
)

type RtspConnectionParams struct {
	IP           string
	Port         int
	Path         string
	Transmission headers.RtspTransmissionType
	Transport    headers.RtspTransportType
	Credentials  auth.Credentials
}

type RtspConnectionContext struct {
	CSeq          int
	AuthHeader    auth.RtspAuthHeader
	SessionHeader headers.SessionHeader
}

type RtspClientIfc interface {
	Connect(connectionParams RtspConnectionParams) error
	Disconnect()
	SendCommand(rtspCommand commands.RtspCommand) (*responses.RtspResponse, error)
	StartListening(transportHeader headers.TransportHeader)
	GetContext() *RtspConnectionContext
}
