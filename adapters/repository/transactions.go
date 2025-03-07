package repository

import (
	"context"
	"database/sql"

	"github.com/go-faster/errors"
	"github.com/uptrace/bun"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.DatabaseTransactionPort = (*TxRepository)(nil)
var _ ports.DatabaseWithinTransactionPort = (*TxRepository)(nil)

var ErrNoTxInContext = errors.New("no transaction in context")

type (
	contextKey struct{}

	TxRepository struct {
		dbConn *bun.DB
	}
)

func NewTxRepository(dbConn *bun.DB) TxRepository {
	return TxRepository{dbConn: dbConn}
}

func (t TxRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := t.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return ctx, errors.Wrap(err, "begin tx")
	}
	return context.WithValue(ctx, contextKey{}, tx), nil
}

func (t TxRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(contextKey{}).(*bun.Tx)
	if !ok {
		return ErrNoTxInContext
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}
	return nil
}

func (t TxRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(contextKey{}).(*bun.Tx)
	if !ok {
		return ErrNoTxInContext
	}

	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, "rollback tx")
	}
	return nil
}

func (t TxRepository) GetTxOrConn(ctx context.Context) bun.IDB {
	if tx, ok := ctx.Value(contextKey{}).(*bun.Tx); ok {
		return tx
	}
	return t.dbConn
}

func (t TxRepository) GetTxDB(ctx context.Context) (*bun.Tx, error) {
	tx, ok := ctx.Value(contextKey{}).(*bun.Tx)
	if !ok {
		return nil, ErrNoTxInContext
	}
	return tx, nil
}

func (t TxRepository) WithInTransactionWithOptions(
	ctx context.Context,
	txFunc func(ctx context.Context) error,
	opts *sql.TxOptions,
) (err error) {
	tx, err := t.dbConn.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer func() {
		var tErr error
		if err != nil {
			tErr = tx.Rollback()
		} else {
			tErr = tx.Commit()
		}

		if tErr != nil && !errors.Is(tErr, sql.ErrTxDone) {
			err = tErr
		}
	}()

	err = txFunc(inject(ctx, tx))
	return err
}

func (t TxRepository) WithInTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return t.WithInTransactionWithOptions(ctx, txFunc, nil)
}

func inject(ctx context.Context, tx bun.IDB) context.Context {
	return context.WithValue(ctx, contextKey{}, tx)
}
