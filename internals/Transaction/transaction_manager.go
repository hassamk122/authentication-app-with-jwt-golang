package transactions

import (
	"context"
	"database/sql"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

type TxManager[T any] struct {
	Db *sql.DB
}

func NewTxManager[T any](db *sql.DB) *TxManager[T] {
	return &TxManager[T]{
		Db: db,
	}
}

func (tm *TxManager[T]) StartTransaction(ctx context.Context, fn func(*store.Queries) (T, error)) (res T, err error) {
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
