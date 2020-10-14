package auth

import (
	"fmt"
	"testing"
)

func TestBasicAuthStringCreation(t *testing.T) {
	sut := basicRtspAuthHeader{
		credentials: Credentials{
			Username: "admin",
			Password: "PassWord",
		},
	}
	expected := "Authorization: Basic YWRtaW46UGFzc1dvcmQ="

	if sut.String("", "") != expected {
		t.Errorf(fmt.Sprintf("Failed to provide basic auth header! Provided:\n%s\nExpected:\n%s", sut.String("", ""), expected))
	}
}

func TestDigestAuthStringCreation(t *testing.T) {
	sut := digestRtspAuthHeader{
		realm: "REALM_ACCC8E0CBE87",
		nonce: "01f980f9Y519728610823565637cfb6b67be21b5245120",
		credentials: Credentials{
			Username: "admin",
			Password: "PassWord",
		},
	}
	testMethod := "DESCRIBE"
	testAddress := "rtsp://192.168.0.10:554/bitrate=1000"
	expected := `Authorization: Digest username="admin", realm="REALM_ACCC8E0CBE87", nonce="01f980f9Y519728610823565637cfb6b67be21b5245120", uri="rtsp://192.168.0.10:554/bitrate=1000", response="596307f242f2ead48af0756fc8211409"`

	if sut.String(testMethod, testAddress) != expected {
		t.Errorf(fmt.Sprintf("Failed to provide digest auth header! Provided:\n%s\nExpected:\n%s", sut.String(testMethod, testAddress), expected))
	}
}
