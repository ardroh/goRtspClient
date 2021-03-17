package handlers

type RtspHandlerFactory struct {
}

func (factory RtspHandlerFactory) Create() rtspHandler {
	optionsHandler := rtspOptionsHandler{}
	describeHandler := rtspDescribeHandler{}
	optionsHandler.SetNext(&describeHandler)
	setupHandler := rtspSetupHandler{}
	describeHandler.SetNext(&setupHandler)
	playHandler := rtspPlayHandler{}
	setupHandler.SetNext(&playHandler)
	return &optionsHandler
}
