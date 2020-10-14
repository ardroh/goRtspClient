package auth

import (
	b64 "encoding/base64"
	"fmt"
)

type RtspAuthHeader interface {
	String() string
}

type Credentials struct {
	Username string
	Password string
}

type basicRtspAuthHeader struct {
	credentials Credentials
}

func (header *basicRtspAuthHeader) String() string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", header.credentials.Username, header.credentials.Password)))
	return fmt.Sprintf("Authorization: Basic %s\n", sEnc)
}

func BuildRtspAuthHeader(request RtspAuthRequest, credentials Credentials) RtspAuthHeader {
	switch request.AuthType {
	case RatNone:
		return nil
	case RatBasic:
		return &basicRtspAuthHeader{
			credentials: credentials,
		}
	case RatDigest:
		return nil
	}
	return nil
}
