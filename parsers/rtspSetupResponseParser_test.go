package parsers_test

import (
	"testing"

	"github.com/ardroh/goRtspClient/parsers"
)

func TestCanParseSessionForAxis(t *testing.T) {
	responseLiteral := `2020/10/12 15:45:46 RTSP/1.0 200 OK
CSeq: 3
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode="PLAY"
Server: GStreamer RTSP server
Session: LZoqIul0P3odd1lb; timeout=60
Date: Sun, 11 Oct 2020 14:46:29 GMT`
	sut := parsers.RtspSetupResponseParser{}

	parsedResponse, err := sut.FromString(responseLiteral)

	if err != nil {
		t.Fatalf("Error during parsing: %s", err)
	}

	if parsedResponse == nil {
		t.Fatal("Can't parse response")
	}

	if parsedResponse.SessionInfo.Timeout != 60 {
		t.Errorf("Can't parse timeout!")
	}

	if parsedResponse.SessionInfo.Id != "LZoqIul0P3odd1lb" {
		t.Errorf("Can't parse session ID!")
	}
}
