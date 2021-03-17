package readers

import "testing"

type indexOfRtpTestData struct {
	testBuffer []byte
	hasRtp     bool
	indexOfRtp int
}

var indexOfRtpTestDataList = []indexOfRtpTestData{
	{
		testBuffer: []byte{
			byte('a'),
			byte('a'),
			byte('a'),
			byte('$'),
			byte(0),
			byte('0'),
			byte('0'),
			byte('a'),
		},
		hasRtp:     true,
		indexOfRtp: 3,
	},
	{
		testBuffer: []byte{
			byte('a'),
			byte('a'),
			byte('a'),
			byte('$'),
			byte('a'),
		},
		hasRtp:     false,
		indexOfRtp: -1,
	},
	{
		testBuffer: []byte{
			byte('a'),
			byte('a'),
			byte('a'),
			byte('$'),
			byte('0'),
		},
		hasRtp:     false,
		indexOfRtp: -1,
	},
	{
		testBuffer: []byte{
			byte('a'),
			byte('a'),
			byte('$'),
			byte(0),
		},
		hasRtp:     true,
		indexOfRtp: 2,
	},
}

func TestCanDeduceIndexOfRtpDataInMixedBuffer(t *testing.T) {
	for i, testData := range indexOfRtpTestDataList {
		t.Logf("Test data no. %d", i)
		hasRtp, rtpIndex := getRtpDataStartIndex(testData.testBuffer)
		if testData.hasRtp {
			if !hasRtp || rtpIndex != testData.indexOfRtp {
				t.Fatalf("Received %d instead of %d", rtpIndex, testData.indexOfRtp)
			}
		} else {
			if hasRtp || rtpIndex == testData.indexOfRtp {
				t.Fatalf("Should not deduce RTP")
			}
		}
	}
}
