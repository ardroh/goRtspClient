package handlers

import (
	"github.com/ardroh/goRtspClient/commands"
)

type rtspPlayHandler struct {
	next rtspHandler
}

func (thisHandler *rtspPlayHandler) SetNext(nextHandler rtspHandler) {
	thisHandler.next = nextHandler
}

func (thisHandler *rtspPlayHandler) Handle(request *RtspConnectRequest) {
	if !request.HasMethod(commands.Play) {
		thisHandler.callNext(request)
		return
	}
	playCmd := commands.RtspPlayCommand{
		SessionID: request.Session.Id,
	}
	response, err := request.RtspClient.SendCommand(playCmd)
	if response == nil || err != nil {
		return
	}
	thisHandler.callNext(request)
}

func (thisHandler rtspPlayHandler) callNext(request *RtspConnectRequest) {
	if thisHandler.next == nil {
		return
	}
	thisHandler.next.Handle(request)
}
