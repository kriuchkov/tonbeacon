package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.DatabaseTransactionPort = (*TxRepository)(nil)

var ErrNoTxInContext = errors.New("no transaction in context")

type contextKey string

const dbTxKey contextKey = "dbTx"

type TxRepository struct {
	dbConn *bun.DB
}

func NewTxRepository(dbConn *bun.DB) TxRepository {
	return TxRepository{dbConn: dbConn}
}

func (t TxRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := t.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return ctx, errors.Wrap(err, "begin tx")
	}

	return context.WithValue(ctx, dbTxKey, tx), nil
}

func (t TxRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(dbTxKey).(*bun.Tx)
	if !ok {
		return ErrNoTxInContext
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (t TxRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(dbTxKey).(*bun.Tx)
	if !ok {
		return ErrNoTxInContext
	}

	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, "rollback tx")
	}

	return nil
}

func (t TxRepository) GetTxOrConn(ctx context.Context) bun.IDB {
	if tx, ok := ctx.Value(dbTxKey).(*bun.Tx); ok {
		return tx
	}
	return t.dbConn
}

func (t TxRepository) GetTxDB(ctx context.Context) (*bun.Tx, error) {
	tx, ok := ctx.Value(dbTxKey).(*bun.Tx)
	if !ok {
		return nil, ErrNoTxInContext
	}
	return tx, nil
}
