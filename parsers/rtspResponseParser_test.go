package parsers_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ardroh/goRtspClient/parsers"
	"github.com/ardroh/goRtspClient/responses"
)

func TestThatCanReadCSeqFromString(t *testing.T) {
	expectedCseq := rand.Int()
	input := fmt.Sprintf(`RTSP/1.0 401 Unauthorized
CSeq: %d
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`, expectedCseq)
	sut := parsers.RtspResponseParser{}

	response, err := sut.Parse(input)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if response.GetCSeq() != expectedCseq {
		t.Error("Received: ", response.GetCSeq())
	}
}

func TestThatCanReadResponseCode(t *testing.T) {
	input := `RTSP/1.0 401 Unauthorized
CSeq: 2
Server: GStreamer RTSP server
Date: Tue, 13 Oct 2020 19:19:19 GMT`
	sut := parsers.RtspResponseParser{}

	response, err := sut.Parse(input)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	statusCode := response.GetStatusCode()

	if statusCode != responses.RtspUnauthorized {
		t.Errorf("Expected %d, Received %d", responses.RtspUnauthorized, statusCode)
	}
}
