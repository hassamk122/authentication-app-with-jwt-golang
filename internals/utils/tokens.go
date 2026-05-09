package utils

import (
	"log"
	"os"

	"github.com/google/uuid"
)

type Tokens struct {
	RefreshToken string
	AccessToken  string
}

func GenerateTokens(sessionID uuid.UUID, userID uuid.UUID) (*Tokens, error) {

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	refreshToken, err := GenerateRefreshToken(sessionID, jwtKey)
	if err != nil {
		return nil, err
	}

	accessToken, err := GenerateAccessToken(userID, sessionID, jwtKey)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully generated tokens (utils layer)")

	return &Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil

}
