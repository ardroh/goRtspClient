package parsers_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ardroh/goRtspClient/parsers"
)

func TestCanParseMediaInformationForAxis(t *testing.T) {
	expectedProtocolVersion := rand.Int()
	sdpLiteral := fmt.Sprintf(`v=%d
o=- 1188340656180883 1 IN IP4 192.168.0.10
s=Session streamed with GStreamer
i=rtsp-server
t=0 12345
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
a=transform:-1.000000,0.000000,0.000000;0.000000,-1.000000,0.000000;0.000000,0.000000,1.000000
m=audio 57000/2 RTP/AVP 98
a=test=test`, expectedProtocolVersion)
	sut := parsers.MediaInformationParser{}

	mediaInfo, err := sut.FromString(sdpLiteral)

	if err != nil {
		t.Fatalf("Can't parse SDP. Error: %s", err)
	}

	if mediaInfo.ProtocolVersion != expectedProtocolVersion {
		t.Errorf("Failed to parse protocol version.")
	}

	sessionOrigin := mediaInfo.SessionOrigin
	if sessionOrigin.Username != "-" ||
		sessionOrigin.SessionId != "1188340656180883" ||
		sessionOrigin.SessionVersion != 1 ||
		sessionOrigin.NetType != "IN" ||
		sessionOrigin.AddrType != "IP4" ||
		sessionOrigin.UnicastAddress != "192.168.0.10" {
		t.Errorf("Can't parse session origin o=!")
	}

	if mediaInfo.SessionName != "Session streamed with GStreamer" {
		t.Errorf("Can't parse session name s=!")
	}

	if *mediaInfo.SessionInfo != "rtsp-server" {
		t.Errorf("Can't parse session info i=!")
	}

	if len(mediaInfo.Timings) != 1 ||
		mediaInfo.Timings[0].StartTime != 0 ||
		mediaInfo.Timings[0].EndTime != 12345 {
		t.Errorf("Can't parse timings!")
	}
	attributes := mediaInfo.SessionAttributes
	if len(attributes) != 4 ||
		!mediaInfo.HasSessionAttribute("tool") ||
		mediaInfo.GetSessionAttribute("tool") != "GStreamer" ||
		!mediaInfo.HasSessionAttribute("type") ||
		mediaInfo.GetSessionAttribute("type") != "broadcast" ||
		!mediaInfo.HasSessionAttribute("range") ||
		mediaInfo.GetSessionAttribute("range") != "npt=now-" ||
		!mediaInfo.HasSessionAttribute("control") ||
		mediaInfo.GetSessionAttribute("control") != "rtsp://192.168.0.10:554/axis-media/media.amp?videocodec=h264&resolution=1280x720&fps=15&videobitrate=1000" {
		t.Error("Can't parse session attributes!")
	}

	medias := mediaInfo.Medias
	if len(medias) != 2 {
		t.Error("Can't parse media tags!")
	}

	media0 := medias[0]
	if media0.Media != "video" ||
		media0.Port.Port != 0 || media0.Port.NumOfPorts != nil ||
		media0.Proto != "RTP/AVP" ||
		media0.Fmt != 96 ||
		media0.Bandwidth == nil || media0.Bandwidth.BwType != "AS" || media0.Bandwidth.Bandwidth != 50000 {
		t.Error("Can't parse media tag props!")
	}

	media0Attr := media0.Attributes
	if len(media0Attr) != 5 {
		t.Error("Can't parse media0 attributes!")
	}

	media1 := medias[1]
	if media1.Media != "audio" ||
		media1.Port.Port != 57000 || *media1.Port.NumOfPorts != 2 ||
		media1.Proto != "RTP/AVP" ||
		media1.Fmt != 98 {
		t.Error("Can't parse media1 tag props!")
	}
}
