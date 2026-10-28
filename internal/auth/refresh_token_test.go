package auth

import (
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	refreshToken, err := MakeRefreshToken()

	if err != nil {
		t.Fatalf("Failed to make RefreshToken: %v", err)
	}

	if len(refreshToken) != 64 {
		t.Fatalf("Expected refresh token length of 64, got %d", len(refreshToken))
	}
}
