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
	tx            ports.DatabaseWithinTransactionPort
	walletManager ports.WalletPort
	database      ports.DatabasePort
	eventManager  ports.OutboxMessagePort
}

func New(options Options) *Account {
	if err := validator.New().Struct(&options); err != nil {
		log.Panic().Err(err).Msg("invalid options")
	}

	options.SetDefaults()

	return &Account{
		walletManager: options.WalletManager,
		tx:            options.TxManager,
		database:      options.DatabaseManager,
		eventManager:  options.EventManager,
	}
}

func (a *Account) CreateAccount(ctx context.Context, accountID string) (*model.Account, error) {
	exists, err := a.database.IsAccountExists(ctx, accountID)
	if err != nil || exists {
		if err != nil {
			return nil, errors.Wrap(err, "check account exists")
		}
		return nil, model.ErrAccountExists
	}

	var account *model.Account
	err = a.tx.WithInTransaction(ctx, func(ctx context.Context) error {
		account, err = a.database.InsertAccount(ctx, accountID)
		if err != nil {
			return errors.Wrap(err, "insert account")
		}

		wallet, err := a.walletManager.CreateWallet(ctx, account.WalletID)
		if err != nil {
			return errors.Wrap(err, "create wallet")
		}

		account.Address = wallet.WalletAddress()

		if err := a.database.UpdateAccount(ctx, account); err != nil {
			return errors.Wrap(err, "update account")
		}

		if err := a.eventManager.Publish(ctx, model.AccountCreated, account); err != nil {
			return errors.Wrap(err, "publish event")
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "create account")
	}
	return account, nil
}

func (a *Account) GetBalance(ctx context.Context, accountID string) (uint64, error) {
	return 0, nil
}

func (a *Account) CloseAccount(ctx context.Context, accountID string) error {
	return a.tx.WithInTransaction(ctx, func(ctx context.Context) error {
		if err := a.database.CloseAccount(ctx, accountID); err != nil {
			return errors.Wrap(err, "close account")
		}

		if err := a.eventManager.Publish(ctx, model.AccountClosed, accountID); err != nil {
			return errors.Wrap(err, "publish event")
		}
		return nil
	})
}

func (a *Account) ListAccounts(ctx context.Context, filter model.ListAccountFilter) ([]model.Account, error) {
	return a.database.ListAccounts(ctx, filter)
}
