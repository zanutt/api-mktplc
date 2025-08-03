package user

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("test_key")

func GenerateJWT(username, userType string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"type":     userType,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
