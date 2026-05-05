package services

import (
	"context"
	"log"

	transaction "github.com/hassamk122/authentication-app-with-jwt-golang/internals/Transaction"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/errs"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/repo"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/utils"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) error
}

type userService struct {
	TransactionManager transaction.TxManager
	UserRepo           repo.UserRepo
}

func NewUserService(tx transaction.TxManager, userRepo repo.UserRepo) *userService {
	return &userService{
		TransactionManager: tx,
		UserRepo:           userRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) error {
	_, err := s.TransactionManager.StartTransaction(ctx,
		func(qtx *store.Queries) error {
			repo := repo.NewUserRepo(qtx)

			log.Println("Checking if email already exists (service layer)")

			_, err := repo.GetEmailByUser(ctx, email)
			if err == nil {
				log.Println("email already exists (service layer)")
				return errs.ErrEmailTaken
			}

			log.Println("Email does not exists hashing password (service layer)")

			hashedPassword, err := utils.HashPassword(password)
			if err != nil {
				log.Println("hashing failed (service layer)")
				return err
			}

			log.Println("hashed password, trying to save user to db (service layer)")

			_, err = repo.CreateUser(ctx, store.CreateUserParams{
				Username: username,
				Email:    email,
				Password: hashedPassword,
			})

			return err
		})

	return err
}
