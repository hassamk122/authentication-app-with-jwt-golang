package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type RefreshTokenClaims struct {
	SessionId uuid.UUID `json:"session_id"`
	jwt.StandardClaims
}

func GenerateRefreshToken(sessionId uuid.UUID, secretKey []byte) (string, error) {
	claims := RefreshTokenClaims{
		SessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
			Issuer:    "auth_server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
