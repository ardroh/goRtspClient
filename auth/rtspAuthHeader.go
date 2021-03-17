package auth

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
)

type RtspAuthHeader interface {
	String(method string, address string, credentials Credentials) string
}

type Credentials struct {
	Username string
	Password string
}

type basicRtspAuthHeader struct {
}

func (header *basicRtspAuthHeader) String(method string, address string, credentials Credentials) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", credentials.Username, credentials.Password)))
	return fmt.Sprintf("Authorization: Basic %s", sEnc)
}

type digestRtspAuthHeader struct {
	realm string
	nonce string
}

func (header *digestRtspAuthHeader) String(method string, address string, credentials Credentials) string {
	hash1 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", credentials.Username, header.realm, credentials.Password)))
	hash2 := md5.Sum([]byte(fmt.Sprintf("%s:%s", method, address)))
	checksum := md5.Sum([]byte(fmt.Sprintf("%x:%s:%x", hash1, header.nonce, hash2)))
	return fmt.Sprintf(`Authorization: Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%x"`, credentials.Username, header.realm, header.nonce, address, checksum)
}

func BuildRtspAuthHeader(authType RtspAuthType, realm *string, nounce *string) RtspAuthHeader {
	switch authType {
	case RatNone:
		return nil
	case RatBasic:
		return &basicRtspAuthHeader{}
	case RatDigest:
		return &digestRtspAuthHeader{
			realm: *realm,
			nonce: *nounce,
		}
	}
	return nil
}
