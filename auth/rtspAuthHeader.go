package auth

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
)

type RtspAuthHeader interface {
	String(method string, address string) string
}

type Credentials struct {
	Username string
	Password string
}

type basicRtspAuthHeader struct {
	credentials Credentials
}

func (header *basicRtspAuthHeader) String(method string, address string) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", header.credentials.Username, header.credentials.Password)))
	return fmt.Sprintf("Authorization: Basic %s", sEnc)
}

type digestRtspAuthHeader struct {
	realm       string
	nonce       string
	credentials Credentials
}

func (header *digestRtspAuthHeader) String(method string, address string) string {
	hash1 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", header.credentials.Username, header.realm, header.credentials.Password)))
	hash2 := md5.Sum([]byte(fmt.Sprintf("%s:%s", method, address)))
	checksum := md5.Sum([]byte(fmt.Sprintf("%x:%s:%x", hash1, header.nonce, hash2)))
	return fmt.Sprintf(`Authorization: Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%x"`, header.credentials.Username, header.realm, header.nonce, address, checksum)
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
		return &digestRtspAuthHeader{
			realm:       request.Realm,
			nonce:       request.Nonce,
			credentials: credentials,
		}
	}
	return nil
}
