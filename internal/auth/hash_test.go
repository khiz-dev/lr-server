package auth

import "testing"

func TestHashAndComparePassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	match, err := ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("Failed to compare password: %v", err)
	}
	if !match {
		t.Fatalf("Password did not match")
	}
}

func TestHashAndComparePassword_Fail(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	wrongPassword := "wrongpassword"
	match, err := ComparePassword(hash, wrongPassword)
	if err != nil {
		t.Fatalf("Failed to compare password: %v", err)
	}
	if match {
		t.Fatalf("Password should not have matched")
	}
}

func TestHashAndComparePassword_EmptyPassword(t *testing.T) {
	password := ""
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	match, err := ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("Failed to compare password: %v", err)
	}
	if !match {
		t.Fatalf("Password did not match")
	}
}