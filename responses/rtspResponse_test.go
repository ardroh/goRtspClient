package responses

import (
	"testing"
)

func TestNotCrashOnGetSeqWhenEmptyOriginalString(t *testing.T) {
	sut := RtspResponse{
		OriginalString: "",
	}

	sut.GetCSeq()
}
