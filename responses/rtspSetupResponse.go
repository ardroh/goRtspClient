package responses

type SessionInfo struct {
	Timeout int
	Id      string
}

type RtspSetupResponse struct {
	RtspResponse RtspResponse
	SessionInfo  SessionInfo
}
