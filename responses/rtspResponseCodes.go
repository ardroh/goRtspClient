package responses

type RtspResponseCodes int

const (
	RtspUnknown      RtspResponseCodes = -1
	RtspOk                             = 200
	RtspUnauthorized                   = 401
)
