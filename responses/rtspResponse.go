package responses

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ardroh/goRtspClient/auth"
)

type RtspResponse struct {
	OriginalString string
}

func (response *RtspResponse) GetCSeq() int {
	if len(response.OriginalString) == 0 {
		return -1
	}
	r, _ := regexp.Compile("CSeq: (.*?)\n")
	matches := r.FindStringSubmatch(response.OriginalString)
	num, err := strconv.Atoi(strings.TrimSpace(matches[1]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return num
}

func (response *RtspResponse) GetStatusCode() RtspResponseCodes {
	r, _ := regexp.Compile("(RTSP/1.0) (.*?) (.*)")
	matches := r.FindStringSubmatch(response.OriginalString)
	if len(matches) < 1 {
		log.Panicln("Can't get status code!")
		return -1
	}
	num, err := strconv.Atoi(strings.TrimSpace(matches[2]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return RtspResponseCodes(num)
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
