package main

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

const testSecret = "test-secret-key"

func makeAccessToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(testSecret))
}

func parseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(testSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func makeRefreshToken() string {
	return uuid.New().String()
}

func TestHashPassword_ProducesHash(t *testing.T) {
	hash, err := hashPassword("secret123")
	if err != nil {
		t.Fatalf("hashPassword: unexpected error: %v", err)
	}
	if hash == "secret123" {
		t.Fatal("hashPassword: hash must not equal plain-text password")
	}
	if !strings.HasPrefix(hash, "$2a$") {
		t.Errorf("hashPassword: expected bcrypt hash prefix, got %q", hash[:4])
	}
}

func TestHashPassword_DifferentHashesSameInput(t *testing.T) {
	h1, _ := hashPassword("secret123")
	h2, _ := hashPassword("secret123")
	if h1 == h2 {
		t.Error("hashPassword: same input should produce different hashes due to bcrypt salting")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	hash, _ := hashPassword("correct-horse")
	if err := checkPassword(hash, "correct-horse"); err != nil {
		t.Errorf("checkPassword: expected match, got: %v", err)
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, _ := hashPassword("correct-horse")
	if err := checkPassword(hash, "wrong-horse"); err == nil {
		t.Error("checkPassword: expected error for wrong password, got nil")
	}
}

func TestCheckPassword_Empty(t *testing.T) {
	hash, _ := hashPassword("non-empty")
	if err := checkPassword(hash, ""); err == nil {
		t.Error("checkPassword: empty password should not match")
	}
}

func TestMakeAccessToken_ContainsUserID(t *testing.T) {
	token, err := makeAccessToken(99)
	if err != nil {
		t.Fatalf("makeAccessToken: %v", err)
	}
	claims, err := parseAccessToken(token)
	if err != nil {
		t.Fatalf("parseAccessToken: %v", err)
	}
	if int(claims["user_id"].(float64)) != 99 {
		t.Errorf("expected user_id=99, got %v", claims["user_id"])
	}
}

func TestParseAccessToken_ExpiredToken(t *testing.T) {
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	signed, _ := expired.SignedString([]byte(testSecret))
	if _, err := parseAccessToken(signed); err == nil {
		t.Error("parseAccessToken: expected error for expired token")
	}
}

func TestParseAccessToken_WrongSecret(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte("wrong-secret"))
	if _, err := parseAccessToken(signed); err == nil {
		t.Error("parseAccessToken: expected error for wrong secret")
	}
}

func TestParseAccessToken_Tampered(t *testing.T) {
	token, _ := makeAccessToken(1)
	if _, err := parseAccessToken(token + "tampered"); err == nil {
		t.Error("parseAccessToken: expected error for tampered token")
	}
}

func TestMakeRefreshToken_IsValidUUID(t *testing.T) {
	rt := makeRefreshToken()
	if _, err := uuid.Parse(rt); err != nil {
		t.Errorf("makeRefreshToken: expected valid UUID, got %q", rt)
	}
}

func TestMakeRefreshToken_IsUnique(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		rt := makeRefreshToken()
		if _, dup := seen[rt]; dup {
			t.Fatalf("makeRefreshToken: duplicate at iteration %d", i)
		}
		seen[rt] = struct{}{}
	}
}

func TestLoginFlow_Success(t *testing.T) {
	hash, _ := hashPassword("hunter2")
	if err := checkPassword(hash, "hunter2"); err != nil {
		t.Fatalf("login: password check failed: %v", err)
	}
	token, err := makeAccessToken(42)
	if err != nil {
		t.Fatalf("login: token generation failed: %v", err)
	}
	claims, err := parseAccessToken(token)
	if err != nil {
		t.Fatalf("login: could not parse token: %v", err)
	}
	if int(claims["user_id"].(float64)) != 42 {
		t.Error("login: wrong user_id in token")
	}
}

func TestLoginFlow_WrongPassword(t *testing.T) {
	hash, _ := hashPassword("hunter2")
	if err := checkPassword(hash, "wrong"); err == nil {
		t.Error("login: expected failure for wrong password")
	}
}

func TestRefreshFlow_NewTokensDiffer(t *testing.T) {
	old := makeRefreshToken()
	new := makeRefreshToken()
	if old == new {
		t.Error("refresh tokens must be unique")
	}
}