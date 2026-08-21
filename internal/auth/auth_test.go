package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if !match {
		t.Errorf("Expected password to match hash")
	}

	wrongMatch, _ := CheckPasswordHash("wrongPassword", hash)
	if wrongMatch {
		t.Errorf("Expected wrong password not to match")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "supersecretjwtkey1234567890123"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected %v, got %v", userID, parsedID)
	}
}
