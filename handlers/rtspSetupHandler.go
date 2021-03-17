package handlers

import (
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/headers"
	"github.com/ardroh/goRtspClient/parsers"
)

type rtspSetupHandler struct {
	next rtspHandler
}

func (thisHandler *rtspSetupHandler) SetNext(nextHandler rtspHandler) {
	thisHandler.next = nextHandler
}

func (thisHandler *rtspSetupHandler) Handle(request *RtspConnectRequest) {
	if !request.HasMethod(commands.Setup) {
		thisHandler.callNext(request)
		return
	}
	var setupCmd commands.RtspSetupCommand
	switch request.Transport {
	case headers.RtpAvpTcp:
		setupCmd = commands.RtspSetupCommand{
			TransportHeader: headers.TransportHeader{
				TransportType:    request.Transport,
				TransmissionType: request.Transmission,
				InterleavedPair: &headers.InterleavedPair{
					RangeMin: 0,
					RangeMax: 1,
				},
			},
		}
	default:
		return
	}
	response, err := request.RtspClient.SendCommand(setupCmd)
	if response == nil || err != nil {
		return
	}
	parsedResponse, err := parsers.RtspSetupResponseParser{}.FromBaseResponse(*response)
	if err != nil {
		return
	}
	request.Session = parsedResponse.SessionHeader
	request.RtspClient.GetContext().SessionHeader = parsedResponse.SessionHeader
	request.RtspClient.StartListening(parsedResponse.Transport)
	thisHandler.callNext(request)
}

func (thisHandler rtspSetupHandler) callNext(request *RtspConnectRequest) {
	if thisHandler.next == nil {
		return
	}
	thisHandler.next.Handle(request)
}
