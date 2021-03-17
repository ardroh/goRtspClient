package headers_test

import (
	"testing"

	"github.com/ardroh/goRtspClient/headers"
)

type SessionHeaderTestData struct {
	InputLiteral string
	Expected     headers.SessionHeader
}

var sessionHeaderTestData []SessionHeaderTestData = []SessionHeaderTestData{
	{
		InputLiteral: "Session: LZoqIul0P3odd1lb; timeout=60",
		Expected: headers.SessionHeader{
			Id:      "LZoqIul0P3odd1lb",
			Timeout: 60,
		},
	},
}

func TestSessionHeaderParsing(t *testing.T) {
	for _, testData := range sessionHeaderTestData {
		sut := headers.SessionHeader{}

		parsed, err := sut.FromString(testData.InputLiteral)

		if err != nil {
			t.Fatalf("Can't parse header! Error: %s", err)
		}

		if parsed.Id != testData.Expected.Id {
			t.Errorf("Failed to parse id!")
		}

		if parsed.Timeout != testData.Expected.Timeout {
			t.Errorf("Failed to parse timeout!")
		}
	}
}
