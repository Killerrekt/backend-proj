package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN = iota
	REFRESH_TOKEN
)

type TokenType int

type TokenPayload struct {
	Email   string
	Role    string
	Version int
}

func CreateToken(exp time.Duration, payload TokenPayload, tokenType TokenType, secretKey string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(exp).Unix(),
		"sub": payload.Email,
	}

	if tokenType == ACCESS_TOKEN {
		claims["role"] = payload.Role
		claims["version"] = payload.Version
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))

	return
}
