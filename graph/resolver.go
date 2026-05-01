package graph

import (
	"fmt"
	"time"

	"auth-server/ent"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ContextKey string

type Resolver struct {
	Client *ent.Client
}

func generateToken(userID int) (string, error) {
	secret := []byte("mysecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(secret)
}

func generateRefreshToken() string {
	return uuid.New().String()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims.(jwt.MapClaims), nil
}