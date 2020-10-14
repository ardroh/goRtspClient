package commands

type RtspCommandType string

const (
	RtspOptions  RtspCommandType = "OPTIONS"
	RtspDescribe                 = "DESCRIBE"
	RtspSetup                    = "SETUP"
	RtspPlay                     = "PLAY"
	RtspTeardown                 = "TEARDOWN"
)

type RtspCommand interface {
	GetCommandType() RtspCommandType
	GetParamsString() string
}
