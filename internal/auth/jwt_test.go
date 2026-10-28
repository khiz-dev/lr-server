package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := time.Minute * 5

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if returnedUserID != userID {
		t.Fatalf("Expected userID %v, got %v", userID, returnedUserID)
	}
}

func TestValidateJWTInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	secret := "testsecret"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatalf("Expected error for invalid token, got nil")
	}
}

func TestValidateJWTExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	expiresIn := -time.Minute * 5 // Token already expired

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("Expected error for expired token, got nil")
	}
}