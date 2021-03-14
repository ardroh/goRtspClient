package responses

import "github.com/ardroh/goRtspClient/headers"

type RtspSetupResponse struct {
	RtspResponse  RtspResponse
	SessionHeader headers.SessionHeader
	Transport     headers.TransportHeader
}
