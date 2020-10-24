package responses

import (
	"github.com/ardroh/goRtspClient/auth"
)

type RtspResponse struct {
	OriginalString string
	CSeq           int
	StatusCode     RtspResponseCodes
	AuthRequest    auth.RtspAuthRequest
}
