package parsers_test

import (
	"testing"

	"github.com/ardroh/goRtspClient/parsers"
	"github.com/ardroh/goRtspClient/responses"
)

func TestThatCanReadBasicDataFromOk(t *testing.T) {
	input := `RTSP/1.0 200 OK
	CSeq: 7
	Content-Type: application/sdp
	Content-Base: rtsp://192.168.0.10:554/axis-media/media.amp/
	Server: GStreamer RTSP server
	Date: Sat, 09 May 2020 01:56:03 GMT
	Content-Length: 710`
	sut := parsers.RtspResponseParser{}

	response, err := sut.Parse(input)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if response.StatusCode != responses.RtspOk {
		t.Errorf("Expected %d, Received %d", responses.RtspOk, response.StatusCode)
	}

	if response.Cseq != 7 {
		t.Error("Can't parse CSeq!")
	}

	if response.ContentType == nil || *response.ContentType != "application/sdp" {
		t.Error("Can't parse ContentType!")
	}

	if response.ContentBase == nil || *response.ContentBase != "rtsp://192.168.0.10:554/axis-media/media.amp/" {
		t.Error("Can't parse ContentBase!")
	}

	if response.Server == nil || *response.Server != "GStreamer RTSP server" {
		t.Error("Can't parse Server!")
	}

	if response.DateTime == nil ||
		response.DateTime.Year() != 2020 ||
		response.DateTime.Month() != 5 ||
		response.DateTime.Day() != 9 ||
		response.DateTime.Hour() != 1 ||
		response.DateTime.Minute() != 56 ||
		response.DateTime.Second() != 3 {
		if response.DateTime == nil {
			t.Error("Can't parse date.")
		} else {
			t.Errorf("Received wrong date: %s", response.DateTime.String())
		}
	}

	if response.ContentLength != 710 {
		t.Error("Can't parse ContentLength!")
	}
}
