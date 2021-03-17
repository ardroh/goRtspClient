package commands

type RtspCommandTypes string

const (
	Options      RtspCommandTypes = "OPTIONS"
	Describe                      = "DESCRIBE"
	GetParameter                  = "GET_PARAMETER"
	Pause                         = "PAUSE"
	Play                          = "PLAY"
	Setup                         = "SETUP"
	SetParameter                  = "SET_PARAMETER"
	Teardown                      = "TEARDOWN"
)
