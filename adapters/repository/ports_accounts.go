package repository

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/model"
)

func (d *DatabaseAdapter) IsAccountExists(ctx context.Context, accountID string) (bool, error) {
	var exists bool

	idb := d.GetTxOrConn(ctx)
	err := idb.NewSelect().ColumnExpr("EXISTS(SELECT 1 FROM accounts WHERE id = ?)", accountID).Scan(ctx, &exists)
	if err != nil {
		return false, errors.Wrap(err, "check account exists")
	}

	return exists, nil
}

func (d *DatabaseAdapter) InsertAccount(ctx context.Context, accountID string) (*model.Account, error) {
	idb := d.GetTxOrConn(ctx)

	// Create account object with just the ID, let the database assign wallet_id via sequence
	account := &Account{ID: accountID}
	log.Debug().Any("account", account).Msg("insert account")

	_, err := idb.NewInsert().Model(account).Returning("*").Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "insert account")
	}
	return account.toModel(), nil
}

func (d *DatabaseAdapter) UpdateAccount(ctx context.Context, mAccount *model.Account) error {
	accountModel := fromModelAccount(mAccount)
	log.Debug().Any("account", accountModel).Msg("update account")

	idb := d.GetTxOrConn(ctx)
	if _, err := idb.NewUpdate().Model(accountModel).Where("id = ?", mAccount.ID).Exec(ctx); err != nil {
		return errors.Wrap(err, "update account")
	}
	return nil
}

func (d *DatabaseAdapter) CloseAccount(ctx context.Context, accountID string) error {
	idb := d.GetTxOrConn(ctx)

	if _, err := idb.NewUpdate().Model((*Account)(nil)).Set("is_closed = ?", true).
		Where("id = ?", accountID).Exec(ctx); err != nil {
		return errors.Wrap(err, "close account")
	}
	return nil
}
