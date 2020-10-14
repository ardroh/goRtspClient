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

	if sut.String() != expected {
		t.Errorf(fmt.Sprintf("Failed to provide basic auth header! Provided:\n%s\nExpected:\n%s", sut.String(), expected))
	}
}
