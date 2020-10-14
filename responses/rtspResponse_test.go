package responses

import (
	"testing"

	"github.com/ardroh/goRtspClient/auth"
)

func TestNotCrashOnGetSeqWhenEmptyOriginalString(t *testing.T) {
	sut := RtspResponse{
		OriginalString: "",
	}

	sut.GetCSeq()
}

func TestDeduceNoneAuthOnGetAuthTypeWhenNoAuthProvided(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`

	sut := RtspResponse{
		OriginalString: input,
	}

	if sut.GetRtspAuthType().AuthType != auth.RatNone {
		t.Errorf("Failed to deduce basic auth type!")
	}
}

func TestDeduceBasicAuthOnGetAuthTypeWhenBasicAuthRequested(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
WWW-Authenticate: Basic realm="AXIS_FDSKFLSSAKSL"
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`

	sut := RtspResponse{
		OriginalString: input,
	}

	if sut.GetRtspAuthType().AuthType != auth.RatBasic {
		t.Errorf("Failed to deduce basic auth type!")
	}

	if sut.GetRtspAuthType().Realm != "AXIS_FDSKFLSSAKSL" {
		t.Errorf("Failed to deduce basic realm!")
	}
}

func TestDeduceDigestAuthOnGetAuthTypeWhenDigestRequestedFirst(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
WWW-Authenticate: Digest realm="AXIS_ACCC8E0CBE87", nonce="01f82f1eY2828699f525c781da2d528b24cd9b5f249fbb", stale=FALSE
WWW-Authenticate: Basic realm="AXIS_ACCC8E0CBE87"
Server: GStreamer RTSP server
Date: Fri, 08 May 2020 01:54:32 GMT`

	sut := RtspResponse{
		OriginalString: input,
	}

	if sut.GetRtspAuthType().AuthType != auth.RatDigest {
		t.Errorf("Failed to deduce digest auth type!")
	}

	if sut.GetRtspAuthType().Realm != "AXIS_ACCC8E0CBE87" {
		t.Errorf("Failed to deduce digest realm!")
	}

	if sut.GetRtspAuthType().Nonce != "01f82f1eY2828699f525c781da2d528b24cd9b5f249fbb" {
		t.Errorf("Failed to deduce digest nounce!")
	}
}
