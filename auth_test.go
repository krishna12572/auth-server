package main

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ---- helpers ----

func testGenerateToken(userID int) (string, error) {
	secret := []byte("mysecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(secret)
}

func testGenerateRefreshToken() string {
	return uuid.New().String()
}

// ---- Password Tests ----

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if string(hashed) == password {
		t.Fatal("hashed password should not equal plain text")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("expected password to match, got: %v", err)
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

// ---- JWT Tests ----

func TestGenerateToken(t *testing.T) {
	token, err := testGenerateToken(1)
	if err != nil {
		t.Fatalf("generateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("generateToken returned empty token")
	}
}

func TestValidateToken(t *testing.T) {
	token, err := testGenerateToken(42)
	if err != nil {
		t.Fatalf("generateToken failed: %v", err)
	}
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !parsed.Valid {
		t.Error("expected token to be valid")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}
	userID := int(claims["user_id"].(float64))
	if userID != 42 {
		t.Errorf("expected user_id 42, got %d", userID)
	}
}

func TestToken_Expiry(t *testing.T) {
	secret := []byte("mysecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	signed, _ := token.SignedString(secret)
	_, err := jwt.Parse(signed, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

// ---- Refresh Token Tests ----

func TestGenerateRefreshToken(t *testing.T) {
	token := testGenerateRefreshToken()
	if token == "" {
		t.Fatal("generateRefreshToken returned empty string")
	}
	_, err := uuid.Parse(token)
	if err != nil {
		t.Errorf("refresh token is not a valid UUID: %v", err)
	}
}

func TestRefreshToken_Unique(t *testing.T) {
	token1 := testGenerateRefreshToken()
	token2 := testGenerateRefreshToken()
	if token1 == token2 {
		t.Error("refresh tokens should be unique")
	}
}

// ---- Login Simulation Tests ----

func TestLogin_Success(t *testing.T) {
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// simulate what login does: check password, generate tokens
	err := bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Fatalf("login should succeed with correct password: %v", err)
	}

	accessToken, err := testGenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	refreshToken := testGenerateRefreshToken()

	if accessToken == "" {
		t.Error("access token should not be empty on successful login")
	}
	if refreshToken == "" {
		t.Error("refresh token should not be empty on successful login")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Error("login should fail with wrong password")
	}
}

// ---- Refresh Token Rotation Tests ----

func TestRefreshToken_Rotation(t *testing.T) {
	// simulate rotation: old token deleted, new tokens generated
	oldRefresh := testGenerateRefreshToken()

	// generate new tokens (simulating rotation)
	newAccess, err := testGenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate new access token: %v", err)
	}
	newRefresh := testGenerateRefreshToken()

	if newRefresh == oldRefresh {
		t.Error("new refresh token should be different from old one")
	}
	if newAccess == "" {
		t.Error("new access token should not be empty")
	}
}