package utils

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Tokens struct {
	RefreshToken string
	AccessToken  string
}

type AccessTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionId uuid.UUID `json:"session_id"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	SessionId uuid.UUID `json:"session_id"`
	jwt.StandardClaims
}

func GenerateTokens(sessionID uuid.UUID, userID uuid.UUID) (*Tokens, error) {

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	refreshToken, err := generateRefreshToken(sessionID, jwtKey)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(userID, sessionID, jwtKey)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully generated tokens (utils layer)")

	return &Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil

}

func generateAccessToken(user_id uuid.UUID, sessionId uuid.UUID, secretKey []byte) (string, error) {

	log.Println("user id ,", user_id)
	claims := AccessTokenClaims{
		UserID:    user_id,
		SessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "auth_server",
		},
	}

	log.Println("access_token claims ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func generateRefreshToken(sessionId uuid.UUID, secretKey []byte) (string, error) {
	claims := RefreshTokenClaims{
		SessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 30).Unix(),
			Issuer:    "auth_server",
		},
	}

	log.Println("refresh_token claims ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
