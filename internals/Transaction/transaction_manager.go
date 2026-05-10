package transactions

import (
	"context"
	"database/sql"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type TxManager struct {
	Db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{Db: db}
}

func StartTransaction[T any](ctx context.Context, tm *TxManager, fn func(*store.Queries) (T, error)) (res T, err error) {
	tx, err := tm.Db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	qtx := store.New(tx)
	res, err = fn(qtx)
	return res, err
}
