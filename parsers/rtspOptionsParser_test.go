package parsers_test

import (
	"testing"

	"github.com/ardroh/goRtspClient/commands"
	"github.com/ardroh/goRtspClient/parsers"
)

func CheckIsCommandInArray(array []commands.RtspCommandTypes, element commands.RtspCommandTypes) bool {
	for _, v := range array {
		if element == v {
			return true
		}
	}
	return false
}

func TestThatCanReadAvailableMethods(t *testing.T) {
	input := `RTSP/1.0 200 OK
	CSeq: 1
	Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, GET_PARAMETER, SET_PARAMETER`
	sut := parsers.RtspOptionsResponseParser{}
	expectedMethods := []commands.RtspCommandTypes{
		commands.Describe,
		commands.Setup,
		commands.Teardown,
		commands.Play,
		commands.Pause,
		commands.GetParameter,
		commands.SetParameter,
	}

	response, err := sut.FromString(input)

	if err != nil {
		t.Errorf("Failed to parse response. Error: %s", err)
	}

	for _, v := range expectedMethods {
		if !CheckIsCommandInArray(response.AvailableMethods, v) {
			t.Errorf("Can't get %s in array.", v)
		}
	}
}
