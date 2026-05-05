package repo

import (
	"context"
	"database/sql"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type UserRepo interface {
	CreateUser(ctx context.Context, arg store.CreateUserParams) (store.CreateUserRow, error)
	GetEmailByUser(ctx context.Context, email string) (store.GetUserByEmailRow, error)
}

type userRepo struct {
	Db      *sql.DB
	queries *store.Queries
}

func NewUserRepo(q *store.Queries) *userRepo {
	return &userRepo{
		queries: q,
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, arg store.CreateUserParams) (store.CreateUserRow, error) {
	return ur.queries.CreateUser(ctx, arg)
}

func (ur *userRepo) GetEmailByUser(ctx context.Context, email string) (store.GetUserByEmailRow, error) {
	return ur.queries.GetUserByEmail(ctx, email)
}
