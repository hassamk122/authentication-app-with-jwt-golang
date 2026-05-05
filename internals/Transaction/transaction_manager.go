package transaction

import (
	"context"
	"database/sql"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type TxManager interface {
	StartTransaction(ctx context.Context, fn func(*store.Queries) error) (*store.Queries, error)
}

type txManager struct {
	Db *sql.DB
}

func NewTxManager(db *sql.DB) *txManager {
	return &txManager{
		Db: db,
	}
}

func (tm *txManager) StartTransaction(ctx context.Context, fn func(*store.Queries) error) (*store.Queries, error) {
	tx, err := tm.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	qtx := store.New(tx)

	err = fn(qtx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return qtx, tx.Commit()
}
