package parsers_test

import (
	"fmt"
	"testing"

	"github.com/ardroh/goRtspClient/parsers"
)

func TestCanParseForAxis(t *testing.T) {
	sdpLiteral := `v=0
o=- 1188340656180883 1 IN IP4 192.168.0.10
s=Session streamed with GStreamer
i=rtsp-server
t=0 0
a=tool:GStreamer
a=type:broadcast
a=range:npt=now-
a=control:rtsp://192.168.0.10:554/axis-media/media.amp?videocodec=h264&resolution=1280x720&fps=15&videobitrate=1000
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:50000
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1;profile-level-id=4d0029;sprop-parameter-sets=Z00AKeKQCgC3YC3AQEBpB4kRUA==,aO48gA==
a=control:rtsp://192.168.0.10:554/axis-media/media.amp/stream=0?videocodec=h264&resolution=1280x720&fps=15&videobitrate=1000
a=framerate:15.000000
a=transform:-1.000000,0.000000,0.000000;0.000000,-1.000000,0.000000;0.000000,0.000000,1.000000`
	responseLiteral := fmt.Sprintf(`RTSP/1.0 200 OK
CSeq: 7
Content-Type: application/sdp
Content-Base: rtsp://192.168.0.10:554/axis-media/media.amp/
Server: GStreamer RTSP server
Date: Sat, 09 May 2020 01:56:03 GMT
Content-Length: 710

%s`, sdpLiteral)
	sut := parsers.RtspDescribeResponseParser{}

	parsedResponse, err := sut.FromString(responseLiteral)

	if err != nil {
		t.Fatalf("Failed to parse. Error: %s", err)
	}

	if parsedResponse.MediaInfo.OriginalLiteral != sdpLiteral {
		t.Errorf("Can't extract SDP! Received: %d Expected: %d", len(parsedResponse.MediaInfo.OriginalLiteral), len(sdpLiteral))
	}
}
