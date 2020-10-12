package responses

import "testing"

func TestCanParseSessionForAxis(t *testing.T) {
	response := `2020/10/12 15:45:46 RTSP/1.0 200 OK
CSeq: 3
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode="PLAY"
Server: GStreamer RTSP server
Session: LZoqIul0P3odd1lb; timeout=60
Date: Sun, 11 Oct 2020 14:46:29 GMT`
	sut := RtspSetupResponse{
		rtspResponse: RtspResponse{
			OriginalString: response,
		},
	}

	if sut.getTimeout() != 60 {
		t.Errorf("Failed to parse timeout!")
	}

	if sut.getSession() != "LZoqIul0P3odd1lb" {
		t.Errorf("Failed to parse session!")
	}
}
