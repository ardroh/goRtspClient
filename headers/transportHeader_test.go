package headers_test

import (
	"testing"

	"github.com/ardroh/goRtspClient/headers"
)

type TransportHeaderTestData struct {
	Literal string
	Object  headers.TransportHeader
}

var transportHeaderTestData []TransportHeaderTestData = []TransportHeaderTestData{
	{
		Literal: "Transport: RTP/AVP/TCP;unicast;interleaved=0-1;ssrc=5F53FB16;mode=\"PLAY\"",
		Object: headers.TransportHeader{
			TransportType:    headers.RtpAvpTcp,
			TransmissionType: headers.Unicast,
			InterleavedPair:  &headers.InterleavedPair{RangeMin: 0, RangeMax: 1},
			Ssrc:             "5F53FB16",
			Mode:             headers.PlayTransportMode,
		},
	},
}

func TestTransportHeaderParsing(t *testing.T) {
	for _, testData := range transportHeaderTestData {
		sut := headers.TransportHeader{}

		parsed, err := sut.FromString(testData.Literal)

		if err != nil {
			t.Fatalf("Can't parse header! Error: %s", err)
		}

		if parsed.TransportType != testData.Object.TransportType {
			t.Errorf("Can't parse Transport Type!")
		}

		if parsed.TransmissionType != testData.Object.TransmissionType {
			t.Errorf("Can't parse transmission type!")
		}

		if *parsed.InterleavedPair != *testData.Object.InterleavedPair {
			t.Errorf("Can't parse interleaved pair!")
		}

		if parsed.Ssrc != testData.Object.Ssrc {
			t.Errorf("Can't parse ssrc!")
		}

		if parsed.Mode != testData.Object.Mode {
			t.Errorf("Can't parse transport mode!")
		}
	}
}

func TestTransportHeaderToString(t *testing.T) {
	for _, testData := range transportHeaderTestData {
		sut := testData.Object

		if sut.ToString() != testData.Literal {
			t.Errorf("Failed to print header. Got: \"%s\" Expected: \"%s\"", sut.ToString(), testData.Literal)
		}
	}
}
