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
	Email string
	Role  string
}

func CreateToken(
	exp time.Duration,
	payload TokenPayload,
	tokenType TokenType,
	secretKey string,
) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(exp).Unix(),
	}

	if tokenType == ACCESS_TOKEN {
		claims["sub"] = payload.Email
		claims["role"] = payload.Role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))

	return
}
