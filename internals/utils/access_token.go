package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionId uuid.UUID `json:"session_id"`
	jwt.StandardClaims
}

func GenerateAccessToken(user_id uuid.UUID, sessionId uuid.UUID, secretKey []byte) (string, error) {
	claims := AccessTokenClaims{
		UserID:    user_id,
		SessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "auth_server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
