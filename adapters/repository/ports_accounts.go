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

	_, err := idb.NewInsert().Model(account).Column("id", "ton_address").Returning("*").Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "insert exec")
	}
	return account.toModel(), nil
}

func (d *DatabaseAdapter) UpdateAccount(ctx context.Context, mAccount *model.Account) error {
	accountModel := fromModelAccount(mAccount)
	log.Debug().Any("account", accountModel).Msg("update account")

	idb := d.GetTxOrConn(ctx)
	if _, err := idb.NewUpdate().Model(accountModel).Where("id = ?", mAccount.ID).Exec(ctx); err != nil {
		return errors.Wrap(err, "update exec")
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

func (d *DatabaseAdapter) ListAccounts(ctx context.Context, filter model.ListAccountFilter) ([]model.Account, error) {
	idb := d.GetTxOrConn(ctx)

	var accounts []Account
	query := idb.NewSelect().Model(&accounts)

	if filter.IsClosed != nil {
		query.Where("is_closed = ?", *filter.IsClosed)
	}

	if filter.WalletIDs != nil {
		query.Where("wallet_id IN (?)", *filter.WalletIDs)
	}

	query.Offset(filter.Offset).Limit(filter.Limit)

	if err := query.Scan(ctx); err != nil {
		return nil, errors.Wrap(err, "list accounts")
	}

	var result []model.Account
	for _, account := range accounts {
		result = append(result, *account.toModel())
	}
	return result, nil
}

func (d *DatabaseAdapter) GetWalletIDByAccountID(ctx context.Context, accountID string) (uint32, error) {
	idb := d.GetTxOrConn(ctx)

	var walletID uint32
	err := idb.NewSelect().ColumnExpr("wallet_id").Model((*Account)(nil)).Where("id = ?", accountID).Scan(ctx, &walletID)
	if err != nil {
		return 0, errors.Wrap(err, "get wallet id by account id")
	}
	return walletID, nil
}
