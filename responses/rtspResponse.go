package responses

import (
	"regexp"
	"strings"
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
}

func (response *RtspResponse) GetCSeq() int {
	return response.Cseq
}

func (response *RtspResponse) GetStatusCode() RtspResponseCodes {
	return response.StatusCode
}

func getRealm(authString string) string {
	r, _ := regexp.Compile("realm=\"(.*?)\"")
	matches := r.FindStringSubmatch(authString)
	return matches[1]
}

func getNonce(authString string) string {
	r, _ := regexp.Compile("nonce=\"(.*?)\"")
	matches := r.FindStringSubmatch(authString)
	return matches[1]
}

func (response *RtspResponse) GetRtspAuthType() auth.RtspAuthRequest {
	r, _ := regexp.Compile("(WWW-Authenticate:) (.*)")
	matches := r.FindStringSubmatch(response.OriginalString)
	if len(matches) < 1 {
		return auth.RtspAuthRequest{
			AuthType: auth.RatNone,
		}
	}
	if strings.Contains(matches[2], "Basic") {
		return auth.RtspAuthRequest{
			AuthType: auth.RatBasic,
			Realm:    getRealm(matches[2]),
		}
	}
	return auth.RtspAuthRequest{
		AuthType: auth.RatDigest,
		Realm:    getRealm(matches[2]),
		Nonce:    getNonce(matches[2]),
	}
}
