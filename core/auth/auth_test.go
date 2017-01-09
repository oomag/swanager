package auth

import "testing"

func TestAuthWithEmptyToken(t *testing.T) {
	if WithToken("") {
		t.Error("Should not auth with empty token")
	}
}

func TestAuthWithCurrentToken(t *testing.T) {
	if !WithToken("token") {
		t.Error("Sould authenticate with any token")
	}
}
