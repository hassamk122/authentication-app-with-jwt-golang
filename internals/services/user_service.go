package services

import (
	"context"
	"log"
	"net/http"
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
	Register(ctx context.Context, username, email, password string) (dtos.RegisterInfo, error)
	Login(ctx context.Context, email, password string) (dtos.LoginInfo, error)
	Logout(ctx context.Context, accessToken *http.Cookie) error
}

type userService struct {
	TxManager            transaction.TxManager
	UserRepo             repo.UserRepo
	verificationCodeRepo repo.VerificationCodeRepo
	userSessionRepo      repo.UserSessionRepo
}

func NewUserService(
	TxManager transaction.TxManager,
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

func (s *userService) Register(ctx context.Context, username, email, password string) (dtos.RegisterInfo, error) {
	registerInfo, err := transaction.StartTransaction(ctx, &s.TxManager,
		func(qtx *store.Queries) (dtos.RegisterInfo, error) {
			userRepo := repo.NewUserRepo(qtx)
			verificationCodeRepo := repo.NewVerificationCodeRepo(qtx)
			userSessionRepo := repo.NewUserSessionRepo(qtx)

			log.Println("Checking if email already exists (service layer)")

			user, err := verifyAndSaveUser(ctx, userRepo, username, email, password)
			if err != nil {
				return dtos.RegisterInfo{}, err
			}

			verificationCode, err := verificationCodeRepo.SaveVerificationCodeUser(ctx, store.SaveVerificationCodeParams{
				UserID:          user.PublicID,
				VerficationType: _types.EmailVerification,
				ExpiresAt:       time.Now().AddDate(1, 0, 0),
			})
			if err != nil {
				return dtos.RegisterInfo{}, err
			}

			log.Println("Verification Code generated (service layer)", verificationCode)

			tokens, err := saveSessionAndGenerateTokens(ctx, user.PublicID, userSessionRepo)
			if err != nil {
				return dtos.RegisterInfo{}, err
			}

			return dtos.RegisterInfo{
				User:         user,
				RefreshToken: tokens.RefreshToken,
				AccessToken:  tokens.AccessToken,
			}, err
		})

	if err != nil {
		return dtos.RegisterInfo{}, err
	}

	return registerInfo, err
}

func (s *userService) Login(ctx context.Context, email, password string) (dtos.LoginInfo, error) {
	LoginInfo, err := transaction.StartTransaction(ctx, &s.TxManager,
		func(qtx *store.Queries) (dtos.LoginInfo, error) {

			userRepo := repo.NewUserRepo(qtx)
			userSessionRepo := repo.NewUserSessionRepo(qtx)

			user, sessionID, err := verifyUserAndCreateSession(ctx, email, password, userRepo, userSessionRepo)
			if err != nil {
				return dtos.LoginInfo{}, err
			}

			log.Println("created session (service layer)")

			tokens, err := utils.GenerateTokens(sessionID, user.PublicID)
			if err != nil {
				return dtos.LoginInfo{}, err
			}

			return dtos.LoginInfo{
				User:         user,
				RefreshToken: tokens.RefreshToken,
				AccessToken:  tokens.AccessToken,
			}, nil
		})
	if err != nil {
		return dtos.LoginInfo{}, err
	}

	return LoginInfo, nil
}

func (s *userService) Logout(ctx context.Context, accessToken *http.Cookie) error {

	claims, err := utils.ParseAccessToken(accessToken.Value, utils.GetSecretKey())
	if err != nil {
		return err
	}

	err = s.userSessionRepo.DeleteUserSession(ctx, claims.SessionId)
	if err != nil {
		return err
	}

	return nil
}

func verifyUserAndCreateSession(ctx context.Context, email, password string, userRepo repo.UserRepo, userSessionRepo repo.UserSessionRepo) (*store.GetUserByEmailRow, uuid.UUID, error) {
	user, err := isAuthenticUser(ctx, email, password, userRepo)
	if err != nil {
		log.Println("email does not exists (service layer)")
		return nil, uuid.Nil, err
	}

	log.Println("valid user (service layer)")

	sessionID, err := createSession(ctx, user.PublicID, userSessionRepo)
	if err != nil {
		log.Println("error creating session (service layer)")
		return nil, uuid.Nil, err
	}

	return user, sessionID, nil
}

func isAuthenticUser(ctx context.Context, email, password string, userRepo repo.UserRepo) (*store.GetUserByEmailRow, error) {
	user, err := findUserByEmail(ctx, email, userRepo)
	if err != nil {
		log.Println("email does not exists (service layer)")
		return nil, err
	}

	log.Println("user found with email (service layer)")

	valid := utils.ComparePassword(user.Password, password)
	if !valid {
		log.Println("invalid user (service layer)")
		return nil, errs.ErrInvalidCredentails
	}

	return user, nil
}

func saveSessionAndGenerateTokens(ctx context.Context, userID uuid.UUID, userSessionRepo repo.UserSessionRepo) (*utils.Tokens, error) {
	sessionID, err := createSession(ctx, userID, userSessionRepo)
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

func createSession(ctx context.Context, userID uuid.UUID, userSessionRepo repo.UserSessionRepo) (uuid.UUID, error) {
	session, err := userSessionRepo.CreateUserSession(ctx, store.CreateUserSessionParams{
		UserID:    userID,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return session.ID, nil
}

func findUserByEmail(ctx context.Context, email string, userRepo repo.UserRepo) (*store.GetUserByEmailRow, error) {
	user, err := userRepo.GetEmailByUser(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func verifyAndSaveUser(ctx context.Context, userRepo repo.UserRepo, username, email, password string) (*store.CreateUserRow, error) {
	_, err := findUserByEmail(ctx, email, userRepo)
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
