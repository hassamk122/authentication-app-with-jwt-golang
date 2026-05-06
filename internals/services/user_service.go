package services

import (
	"context"
	"log"

	transaction "github.com/hassamk122/authentication-app-with-jwt-golang/internals/Transaction"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/_types"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/errs"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/repo"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/utils"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) error
}

type userService struct {
	TxManager            transaction.TxManager[any]
	UserRepo             repo.UserRepo
	verificationCodeRepo repo.VerificationCodeRepo
}

func NewUserService(TxManager transaction.TxManager[any], userRepo repo.UserRepo, verificationCodeRepo repo.VerificationCodeRepo) *userService {
	return &userService{
		TxManager:            TxManager,
		UserRepo:             userRepo,
		verificationCodeRepo: verificationCodeRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) error {
	_, err := s.TxManager.StartTransaction(ctx,
		func(qtx *store.Queries) (any, error) {
			userRepo := repo.NewUserRepo(qtx)
			verificationCodeRepo := repo.NewVerificationCodeRepo(qtx)

			log.Println("Checking if email already exists (service layer)")

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

			verificationCode, err := verificationCodeRepo.SaveVerificationCodeUser(ctx, store.SaveVerificationCodeParams{
				UserID:          int32(user.PublicID.ID()),
				VerficationType: _types.EmailVerification,
				// ExpiresAt:       "",
			})

			log.Println(verificationCode)

			return user, err
		})

	return err
}
