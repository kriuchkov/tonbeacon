package account

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.AccountServicePort = (*Account)(nil)

type Account struct {
	tx            ports.DatabaseTransactionPort
	walletManager ports.WalletPort
	database      ports.DatabasePort
	eventManager  ports.OutboxMessagePort
}

func New(options Options) *Account {
	if err := validator.New().Struct(&options); err != nil {
		log.Panic().Err(err).Msg("invalid options")
	}

	return &Account{
		walletManager: options.WalletManager,
		tx:            options.TxManager,
		database:      options.DatabaseManager,
		eventManager:  options.EventManager,
	}
}

func (a *Account) CreateAccount(ctx context.Context, accountID string) (*model.Account, error) {
	if exists, err := a.database.IsAccountExists(ctx, accountID); err != nil || exists {
		if err != nil {
			return nil, errors.Wrap(err, "check account exists")
		}
		return nil, model.ErrAccountExists
	}

	ctx, err := a.tx.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "begin transaction")
	}
	defer func() { a.tx.Rollback(ctx) }() //nolint:errcheck // we don't care about rollback errors

	account, err := a.database.InsertAccount(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "insert account")
	}

	wallet, err := a.walletManager.CreateWallet(ctx, account.WalletID)
	if err != nil {
		return nil, errors.Wrap(err, "create wallet")
	}

	account.Address = wallet.WalletAddress()

	if err := a.database.UpdateAccount(ctx, account); err != nil {
		return nil, errors.Wrap(err, "update account")
	}

	if err := a.eventManager.Publish(ctx, model.AccountCreated, account); err != nil {
		return nil, errors.Wrap(err, "publish event")
	}

	if err := a.tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "commit transaction")
	}
	return account, nil
}

func (a *Account) GetBalance(ctx context.Context, accountID string) (uint64, error) {
	return 0, nil
}

func (a *Account) CloseAccount(ctx context.Context, accountID string) error {
	ctx, err := a.tx.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer a.tx.Rollback(ctx) //nolint:errcheck // we don't care about rollback errors

	if err := a.database.CloseAccount(ctx, accountID); err != nil {
		return errors.Wrap(err, "close account")
	}

	if err := a.eventManager.Publish(ctx, model.AccountClosed, accountID); err != nil {
		return errors.Wrap(err, "publish event")
	}

	if err := a.tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}
