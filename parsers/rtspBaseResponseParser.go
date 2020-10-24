package parsers

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ardroh/goRtspClient/auth"
	"github.com/ardroh/goRtspClient/responses"
)

type RtspBaseResponseParser struct {
	literalResponse string
}

func (parser *RtspBaseResponseParser) Parse(literalResponse string) (*responses.RtspResponse, error) {
	parser.literalResponse = literalResponse
	response := responses.RtspResponse{
		OriginalString: literalResponse,
	}
	response.CSeq = parser.parseCSeq()
	if response.CSeq == -1 {
		return nil, errors.New("failed to parse cseq")
	}
	response.StatusCode = parser.parseStatusCode()
	if response.StatusCode == -1 {
		return nil, errors.New("failed to parse status code")
	}
	response.AuthRequest = parser.parseAuthRequest()
	return &response, nil
}

func (parser *RtspBaseResponseParser) parseCSeq() int {
	if len(parser.literalResponse) == 0 {
		return -1
	}
	r, _ := regexp.Compile("CSeq: (.*?)\n")
	matches := r.FindStringSubmatch(parser.literalResponse)
	num, err := strconv.Atoi(strings.TrimSpace(matches[1]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return num
}

func (parser *RtspBaseResponseParser) parseStatusCode() responses.RtspResponseCodes {
	r, _ := regexp.Compile("(RTSP/1.0) (.*?) (.*)")
	matches := r.FindStringSubmatch(parser.literalResponse)
	if len(matches) < 1 {
		log.Panicln("Can't get status code!")
		return -1
	}
	num, err := strconv.Atoi(strings.TrimSpace(matches[2]))
	if err != nil {
		log.Panicln(err)
		return -1
	}
	return responses.RtspResponseCodes(num)
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

func (parser *RtspBaseResponseParser) parseAuthRequest() auth.RtspAuthRequest {
	r, _ := regexp.Compile("(WWW-Authenticate:) (.*)")
	matches := r.FindStringSubmatch(parser.literalResponse)
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
