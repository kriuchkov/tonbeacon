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

	_, err := idb.NewInsert().Model(transactionModel).Returning("id").Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "insert exec")
	}
	return transactionModel.toModel(), nil
}

func (d *DatabaseAdapter) GetTransactions(ctx context.Context, limit, offset int) ([]*model.Transaction, error) {
	idb := d.GetTxOrConn(ctx)

	transactions := make([]*Transaction, 0)
	err := idb.NewSelect().Model(&transactions).Limit(limit).Offset(offset).Order("lt DESC").Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "select scan")
	}

	result := make([]*model.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		result = append(result, transaction.toModel())
	}
	return result, nil
}
