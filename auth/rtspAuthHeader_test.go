package auth

import (
	"fmt"
	"testing"
)

func TestBasicAuthStringCreation(t *testing.T) {
	credentials := Credentials{
		Username: "admin",
		Password: "PassWord",
	}
	sut := basicRtspAuthHeader{}
	expected := "Authorization: Basic YWRtaW46UGFzc1dvcmQ="

	receivedAuth := sut.String("", "", credentials)

	if receivedAuth != expected {
		t.Errorf(fmt.Sprintf("Failed to provide basic auth header! Provided:\n%s\nExpected:\n%s", receivedAuth, expected))
	}
}

func TestDigestAuthStringCreation(t *testing.T) {
	credentials := Credentials{
		Username: "admin",
		Password: "PassWord",
	}
	sut := digestRtspAuthHeader{
		realm: "REALM_ACCC8E0CBE87",
		nonce: "01f980f9Y519728610823565637cfb6b67be21b5245120",
	}
	testMethod := "DESCRIBE"
	testAddress := "rtsp://192.168.0.10:554/bitrate=1000"
	expected := `Authorization: Digest username="admin", realm="REALM_ACCC8E0CBE87", nonce="01f980f9Y519728610823565637cfb6b67be21b5245120", uri="rtsp://192.168.0.10:554/bitrate=1000", response="596307f242f2ead48af0756fc8211409"`

	receivedAuth := sut.String(testMethod, testAddress, credentials)

	if receivedAuth != expected {
		t.Errorf(fmt.Sprintf("Failed to provide digest auth header! Provided:\n%s\nExpected:\n%s", receivedAuth, expected))
	}
}
