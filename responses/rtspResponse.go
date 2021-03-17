package responses

import (
	"time"

	"github.com/ardroh/goRtspClient/auth"
)

type RtspResponse struct {
	OriginalString string
	StatusCode     RtspResponseCodes
	Cseq           int
	ContentType    *string
	ContentBase    *string
	Server         *string
	DateTime       *time.Time
	ContentLength  int
	AuthHeaders    []auth.RtspAuthHeader
}

func (response *RtspResponse) GetCSeq() int {
	return response.Cseq
}

func (response *RtspResponse) GetStatusCode() RtspResponseCodes {
	return response.StatusCode
}
