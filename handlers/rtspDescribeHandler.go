package handlers

import (
	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/parsers"
)

type rtspDescribeHandler struct {
	next rtspHandler
}

func (thisHandler rtspDescribeHandler) SetNext(nextHandler rtspHandler) {
	thisHandler.next = nextHandler
}

func (thisHandler rtspDescribeHandler) Handle(request *RtspConnectRequest) {
	if !request.HasMethod(commands.Describe) {
		thisHandler.callNext(request)
		return
	}
	describeCmd := commands.RtspDescribeCommand{}
	response, err := request.RtspClient.SendCommand(describeCmd)
	if response == nil || err != nil {
		return
	}
	parsedDescribe, err := parsers.RtspDescribeResponseParser{}.FromBaseResponse(*response)
	if parsedDescribe == nil || err != nil {
		return
	}
	request.Sdp = parsedDescribe.Sdp
	thisHandler.callNext(request)
}

func (thisHandler rtspDescribeHandler) callNext(request *RtspConnectRequest) {
	if thisHandler.next == nil {
		return
	}
	thisHandler.next.Handle(request)
}
