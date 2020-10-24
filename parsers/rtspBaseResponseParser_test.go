package parsers

import (
	"testing"

	"github.com/ardroh/goRtspClient/auth"
	"github.com/ardroh/goRtspClient/responses"
)

func TestBaseResponseParsing(t *testing.T) {
	input := `2020/10/12 15:45:46 RTSP/1.0 200 OK
CSeq: 3
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode="PLAY"
Server: GStreamer RTSP server
Session: LZoqIul0P3odd1lb; timeout=60
Date: Sun, 11 Oct 2020 14:46:29 GMT`
	sut := RtspBaseResponseParser{}

	response, err := sut.Parse(input)

	if err != nil {
		t.Errorf("Error different then null %+v", err)
	}

	if response == nil {
		t.Error("Failed to parse response but no error")
	}

	if response.OriginalString != input {
		t.Error("OriginalString not set")
	}

	if response.CSeq != 3 {
		t.Error("CSeq not parsed")
	}

	if response.StatusCode != responses.RtspOk {
		t.Error("Status code not parsed")
	}
}

func TestDeduceNoneAuthOnGetAuthTypeWhenNoAuthProvided(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`
	sut := RtspBaseResponseParser{}

	response, err := sut.Parse(input)

	if response == nil || err != nil {
		t.Error("Failed to parse response or auth request")
		return
	}

	if response.AuthRequest.AuthType != auth.RatNone {
		t.Errorf("Failed to deduce basic auth type!")
	}
}

func TestDeduceBasicAuth(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
WWW-Authenticate: Basic realm="AXIS_FDSKFLSSAKSL"
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`
	sut := RtspBaseResponseParser{}

	response, err := sut.Parse(input)

	if response == nil || err != nil {
		t.Error("Failed to parse response or auth request")
		return
	}

	if response.AuthRequest.AuthType != auth.RatBasic {
		t.Errorf("Failed to deduce basic auth type!")
		return
	}

	if response.AuthRequest.Realm != "AXIS_FDSKFLSSAKSL" {
		t.Errorf("Failed to deduce basic realm!")
		return
	}
}

func TestDeduceDigestAuthWhenDigestRequestedFirst(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
WWW-Authenticate: Digest realm="AXIS_ACCC8E0CBE87", nonce="01f82f1eY2828699f525c781da2d528b24cd9b5f249fbb", stale=FALSE
WWW-Authenticate: Basic realm="AXIS_ACCC8E0CBE87"
Server: GStreamer RTSP server
Date: Fri, 08 May 2020 01:54:32 GMT`
	sut := RtspBaseResponseParser{}

	response, err := sut.Parse(input)

	if response == nil || err != nil {
		t.Error("Failed to parse response or auth request")
		return
	}

	if response.AuthRequest.AuthType != auth.RatDigest {
		t.Errorf("Failed to deduce digest auth type!")
	}

	if response.AuthRequest.Realm != "AXIS_ACCC8E0CBE87" {
		t.Errorf("Failed to deduce digest realm!")
	}

	if response.AuthRequest.Nonce != "01f82f1eY2828699f525c781da2d528b24cd9b5f249fbb" {
		t.Errorf("Failed to deduce digest nounce!")
	}
}

func TestDeduceBasicAuthWhenBasicRequestedFirst(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
WWW-Authenticate: Basic realm="AXIS_ACCC8E0CBE87"
WWW-Authenticate: Digest realm="AXIS_ACCC8E0CBE87", nonce="01f82f1eY2828699f525c781da2d528b24cd9b5f249fbb", stale=FALSE
Server: GStreamer RTSP server
Date: Fri, 08 May 2020 01:54:32 GMT`
	sut := RtspBaseResponseParser{}

	response, err := sut.Parse(input)

	if response == nil || err != nil {
		t.Error("Failed to parse response or auth request")
		return
	}

	if response.AuthRequest.AuthType != auth.RatBasic {
		t.Errorf("Failed to deduce basic auth type!")
		return
	}

	if response.AuthRequest.Realm != "AXIS_ACCC8E0CBE87" {
		t.Errorf("Failed to deduce basic realm!")
		return
	}
}
