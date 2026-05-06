package repo

import (
	"context"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type verificationCodeRepo struct {
	queries *store.Queries
}

type VerificationCodeRepo interface {
	SaveVerificationCodeUser(ctx context.Context, arg store.SaveVerificationCodeParams) (store.VerificationCode, error)
}

func NewVerificationCodeRepo(q *store.Queries) *verificationCodeRepo {
	return &verificationCodeRepo{
		queries: q,
	}
}

func (ur *verificationCodeRepo) SaveVerificationCodeUser(ctx context.Context, arg store.SaveVerificationCodeParams) (store.VerificationCode, error) {
	return ur.queries.SaveVerificationCode(ctx, arg)
}
