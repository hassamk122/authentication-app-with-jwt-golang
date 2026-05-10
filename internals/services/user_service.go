package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	transaction "github.com/hassamk122/authentication-app-with-jwt-golang/internals/Transaction"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/_types"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/dtos"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/errs"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/repo"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/utils"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (any, error)
}

type userService struct {
	TxManager            transaction.TxManager[any]
	UserRepo             repo.UserRepo
	verificationCodeRepo repo.VerificationCodeRepo
	userSessionRepo      repo.UserSessionRepo
}

func NewUserService(
	TxManager transaction.TxManager[any],
	userRepo repo.UserRepo,
	verificationCodeRepo repo.VerificationCodeRepo,
	userSessionRepo repo.UserSessionRepo) *userService {
	return &userService{
		TxManager:            TxManager,
		UserRepo:             userRepo,
		verificationCodeRepo: verificationCodeRepo,
		userSessionRepo:      userSessionRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (any, error) {
	registerInfo, err := s.TxManager.StartTransaction(ctx,
		func(qtx *store.Queries) (any, error) {
			userRepo := repo.NewUserRepo(qtx)
			verificationCodeRepo := repo.NewVerificationCodeRepo(qtx)
			userSessionRepo := repo.NewUserSessionRepo(qtx)

			log.Println("Checking if email already exists (service layer)")

			user, err := verifyAndSaveUser(ctx, userRepo, username, email, password)
			if err != nil {
				return nil, err
			}

			verificationCode, err := verificationCodeRepo.SaveVerificationCodeUser(ctx, store.SaveVerificationCodeParams{
				UserID:          user.PublicID,
				VerficationType: _types.EmailVerification,
				ExpiresAt:       time.Now().AddDate(1, 0, 0),
			})
			if err != nil {
				return nil, err
			}

			log.Println("Verification Code generated (service layer)", verificationCode)

			tokens, err := SaveSessionAndGenerateTokens(ctx, user.PublicID, userSessionRepo)
			if err != nil {
				return nil, err
			}

			return dtos.RegisterInfo{
				User:         user,
				RefreshToken: tokens.RefreshToken,
				AccessToken:  tokens.AccessToken,
			}, err
		})

	if err != nil {
		return nil, err
	}

	return registerInfo, err
}

func SaveSessionAndGenerateTokens(ctx context.Context, userID uuid.UUID, userSessionRepo repo.UserSessionRepo) (*utils.Tokens, error) {
	sessionID, err := CreateSession(ctx, userID, userSessionRepo)
	if err != nil {
		return nil, err
	}

	log.Println("Session added to db (service layer)")

	tokens, err := utils.GenerateTokens(sessionID, userID)
	if err != nil {
		return nil, err
	}

	return tokens, err
}

func CreateSession(ctx context.Context, userID uuid.UUID, userSessionRepo repo.UserSessionRepo) (uuid.UUID, error) {
	session, err := userSessionRepo.CreateUserSession(ctx, store.CreateUserSessionParams{
		UserID:    userID,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return session.ID, nil
}

func verifyAndSaveUser(ctx context.Context, userRepo repo.UserRepo, username, email, password string) (*store.CreateUserRow, error) {
	_, err := userRepo.GetEmailByUser(ctx, email)
	if err == nil {
		log.Println("email already exists (service layer)")
		return nil, errs.ErrEmailTaken
	}

	log.Println("Email does not exists hashing password (service layer)")

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("hashing failed (service layer)")
		return nil, err
	}

	log.Println("hashed password, trying to save user to db (service layer)")

	user, err := userRepo.CreateUser(ctx, store.CreateUserParams{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	return &user, nil

}
