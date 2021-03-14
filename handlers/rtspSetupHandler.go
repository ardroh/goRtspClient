package handlers

import (
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/parsers"
)

type rtspSetupHandler struct {
	next rtspHandler
}

func (thisHandler rtspSetupHandler) SetNext(nextHandler rtspHandler) {
	thisHandler.next = nextHandler
}

func (thisHandler rtspSetupHandler) Handle(request *RtspConnectRequest) {
	if !request.HasMethod(commands.Describe) {
		thisHandler.callNext(request)
		return
	}
	var setupCmd commands.RtspSetupCommand
	switch request.Transport {
	case commands.RtpAvpTcp:
		setupCmd = commands.RtspSetupCommand{
			Transport:    request.Transport,
			Transmission: request.Transmission,
			InterleavedPair: commands.InterleavedPair{
				RangeMin: 0,
				RangeMax: 1,
			},
		}
	default:
		return
	}
	response, err := request.RtspClient.SendCommand(setupCmd)
	if response == nil || err != nil {
		return
	}
	parsedSetup, err := parsers.RtspSetupResponseParser{}.FromBaseResponse(*response)
	if parsedSetup == nil || err != nil {
		return
	}
	thisHandler.callNext(request)
}

func (thisHandler rtspSetupHandler) callNext(request *RtspConnectRequest) {
	if thisHandler.next == nil {
		return
	}
	thisHandler.next.Handle(request)
}
