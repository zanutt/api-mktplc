package user

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("test_key")

func GenerateJWT(userID uint, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    userType,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
