package services

import (
	transaction "github.com/hassamk122/authentication-app-with-jwt-golang/internals/Transaction"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/repo"
)

type UserService interface {
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
