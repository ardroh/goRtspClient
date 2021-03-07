package handlers

type rtspHandler interface {
	SetNext(handler *rtspHandler)
	Handle(request *RtspConnectRequest)
}
