package auth

import "testing"

func TestAuthWithEmptyToken(t *testing.T) {
	user, err := WithToken("")
	if err == nil || user != nil {
		t.Error("Should not auth with empty token")
	}
}

func TestAuthWithCurrentToken(t *testing.T) {
	user, err := WithToken("token")
	if err != nil || user == nil {
		t.Error("Sould authenticate with any token")
	}
}
