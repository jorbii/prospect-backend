package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		secretKey: secretKey,
	}
}

func (t *TokenService) GenerateToken(userID string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.secretKey))
}
