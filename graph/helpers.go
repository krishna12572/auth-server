package graph

import (
"time"

"github.com/golang-jwt/jwt/v5"
"github.com/google/uuid"
)

var jwtSecret = []byte("mysecret")

func generateToken(userID int) (string, error) {
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
"user_id": userID,
"exp":     time.Now().Add(time.Hour).Unix(),
})
return token.SignedString(jwtSecret)
}

func generateRefreshToken() string {
return uuid.New().String()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
return jwtSecret, nil
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
