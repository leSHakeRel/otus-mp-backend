package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashed == password {
		t.Error("Hashed password should not be equal to original password")
	}

	if len(hashed) == 0 {
		t.Error("Hashed password should not be empty")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	match := CheckPassword(password, hashed)
	if !match {
		t.Error("Correct password should match")
	}

	// Test wrong password
	wrongPassword := "wrongpassword"
	match = CheckPassword(wrongPassword, hashed)
	if match {
		t.Error("Wrong password should not match")
	}
}

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24 * time.Hour
	userID := uuid.New()
	email := "test@example.com"

	jwtService := NewJWTService(secret, expiration)

	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24 * time.Hour
	userID := uuid.New()
	email := "test@example.com"

	jwtService := NewJWTService(secret, expiration)

	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected userID to be %s, got %s", userID.String(), claims.UserID.String())
	}

	if claims.Email != email {
		t.Errorf("Expected email to be '%s', got '%s'", email, claims.Email)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 24 * time.Hour

	jwtService := NewJWTService(secret, expiration)
	invalidToken := "invalid-token"

	_, err := jwtService.ValidateToken(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := 1 * time.Millisecond
	userID := uuid.New()
	email := "test@example.com"

	jwtService := NewJWTService(secret, expiration)

	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = jwtService.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}
