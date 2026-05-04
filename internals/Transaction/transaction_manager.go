package transaction

import (
	"context"
	"database/sql"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(*store.Queries) error) error
}

type txManager struct {
	Db *sql.DB
}

func NewTxManager(db *sql.DB) *txManager {
	return &txManager{
		Db: db,
	}
}

func (tm *txManager) WithTx(ctx context.Context, fn func(*store.Queries) error) error {
	tx, err := tm.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := store.New(tx)

	err = fn(qtx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
