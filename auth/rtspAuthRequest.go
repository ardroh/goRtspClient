package auth

type RtspAuthType int

const (
	RatNone   RtspAuthType = 0
	RatBasic               = 1
	RatDigest              = 2
)

type RtspAuthRequest struct {
	AuthType RtspAuthType
	Realm    string
	Nonce    string
	Address  string
}
