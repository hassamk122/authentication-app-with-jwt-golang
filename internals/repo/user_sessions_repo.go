package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type UserSessionRepo interface {
	CreateUserSession(ctx context.Context, userID uuid.UUID) (store.UserSession, error)
}

type userSessionsRepo struct {
	queries *store.Queries
}

func NewUserSessionRepo(q *store.Queries) *userSessionsRepo {
	return &userSessionsRepo{
		queries: q,
	}
}

func (usr *userSessionsRepo) CreateUserSession(ctx context.Context, userID uuid.UUID) (store.UserSession, error) {
	return usr.queries.CreateUserSession(ctx, userID)
}
