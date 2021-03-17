package handlers

import (
	"log"

	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/parsers"
)

type rtspOptionsHandler struct {
	next rtspHandler
}

func (thisHandler rtspOptionsHandler) SetNext(nextHandler rtspHandler) {
	thisHandler.next = nextHandler
}

func (thisHandler rtspOptionsHandler) Handle(request *RtspConnectRequest) {
	optionsCmd := commands.RtspOptionsCommand{}
	response, err := request.RtspClient.SendCommand(optionsCmd)
	if response == nil || err != nil {
		return
	}
	optionsResponse, err := parsers.RtspOptionsResponseParser{}.FromBaseResponse(*response)
	if err != nil {
		log.Printf("rtspOptionsHandler> Failed to handle RTSP/OPTIONS response! Error: %s", err)
		return
	}
	request.AvailableMethods = optionsResponse.AvailableMethods
	thisHandler.handleNext(request)
}

func (thisHandler rtspOptionsHandler) handleNext(request *RtspConnectRequest) {
	if thisHandler.next == nil {
		return
	}
	thisHandler.next.Handle(request)
}
