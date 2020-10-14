package commands

type RtspCommandTypes string

const (
	Options      RtspCommandTypes = "OPTIONS"
	Describe                      = "DESCRIBE"
	GetParameter                  = "GET_PARAMETER"
	Pause                         = "PAUSE"
	Play                          = "Play"
	Setup                         = "Setup"
	SetParameter                  = "SET_PARAMETER"
	Teardown                      = "TEARDOWN"
)
