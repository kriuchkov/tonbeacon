package repository

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/model"
)

func (d *DatabaseAdapter) InsertTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	idb := d.GetTxOrConn(ctx)

	transactionModel := fromModelTransaction(transaction)
	log.Debug().Any("transaction", transactionModel).Msg("insert transaction")

	_, err := idb.NewInsert().Model(transactionModel).Returning("*").Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "insert exec")
	}

	return transactionModel.toModel(), nil
}
