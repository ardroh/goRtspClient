package responses

import "testing"

func TestCanParseForAxis(t *testing.T) {
	response := `RTSP/1.0 200 OK
	CSeq: 7
	Content-Type: application/sdp
	Content-Base: rtsp://192.168.0.10:554/axis-media/media.amp/
	Server: GStreamer RTSP server
	Date: Sat, 09 May 2020 01:56:03 GMT
	Content-Length: 710
	
	v=0
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
	sut := RtspDescribeResponse{
		rtspResponse: RtspResponse{
			OriginalString: response,
		},
	}

	receivedContentBase := sut.GetContentBase()
	if receivedContentBase == nil || *receivedContentBase != "rtsp://192.168.0.10:554/axis-media/media.amp/" {
		t.Errorf("Failed to parse Content-Base!")
	}

	controlUris := sut.GetControlUris()
	if len(controlUris) == 0 || controlUris[0] != "rtsp://192.168.0.10:554/axis-media/media.amp?videocodec=h264&resolution=1280x720&fps=15&videobitrate=1000" || controlUris[1] != "rtsp://192.168.0.10:554/axis-media/media.amp/stream=0?videocodec=h264&resolution=1280x720&fps=15&videobitrate=1000" {
		t.Errorf("Failed to parse control uris!")
	}
}
