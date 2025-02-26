package repository

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/kriuchkov/tonbeacon/internal/core/model"
	"github.com/rs/zerolog/log"
)

func (d *DatabaseAdapter) IsAccountExists(ctx context.Context, accountID string) (bool, error) {
	var exists bool
	err := d.db.NewSelect().ColumnExpr("EXISTS(SELECT 1 FROM accounts WHERE id = ?)", accountID).Scan(ctx, &exists)
	if err != nil {
		return false, errors.Wrap(err, "check account exists")
	}

	return exists, nil
}

func (d *DatabaseAdapter) InsertAccount(ctx context.Context, accountID string) (*model.Account, error) {
	var count int
	err := d.db.NewSelect().Model((*Account)(nil)).ColumnExpr("count(*)").Scan(ctx, &count)
	if err != nil {
		return nil, errors.Wrap(err, "count accounts")
	}

	account := &Account{ID: accountID, SubwalletID: uint32(count + 1)}
	log.Debug().Any("account", account).Msg("insert account")

	_, err = d.db.NewInsert().Model(account).Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "insert account")
	}
	return account.toModel(), nil
}

func (d *DatabaseAdapter) UpdateAccount(ctx context.Context, mAccount *model.Account) error {
	accountModel := fromModelAccount(mAccount)
	log.Debug().Any("account", accountModel).Msg("update account")

	if _, err := d.db.NewUpdate().Model(accountModel).Where("id = ?", mAccount.ID).Exec(ctx); err != nil {
		return errors.Wrap(err, "update account")
	}
	return nil
}

func (d *DatabaseAdapter) CloseAccount(ctx context.Context, accountID string) error {
	if _, err := d.db.NewUpdate().Model((*Account)(nil)).Set("is_closed = ?", true).
		Where("id = ?", accountID).Exec(ctx); err != nil {
		return errors.Wrap(err, "close account")
	}
	return nil
}
