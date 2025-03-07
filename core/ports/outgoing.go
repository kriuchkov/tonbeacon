package ports

import (
	"context"

	"github.com/kriuchkov/tonbeacon/core/model"
)

type WalletPort interface {
	CreateWallet(ctx context.Context, walletID uint32) (model.WalletWrapper, error)
	GetExtraCurrenciesBalance(ctx context.Context, walletID uint32) ([]model.Balance, error)
	GetBalance(ctx context.Context, walletID uint32) (uint64, error)
}

type DatabaseWithinTransactionPort interface {
	WithInTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

type DatabaseTransactionPort interface {
	DatabaseWithinTransactionPort

	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type OutboxMessagePort interface {
	Publish(ctx context.Context, eventType model.EventType, payload any) error
}

type PublisherPort interface {
	Publish(ctx context.Context, message any) error
	Close() error
}

type (
	AccountDatabasePort interface {
		IsAccountExists(ctx context.Context, accountID string) (bool, error)
		InsertAccount(ctx context.Context, accountID string) (*model.Account, error)
		UpdateAccount(ctx context.Context, account *model.Account) error
		CloseAccount(ctx context.Context, accountID string) error
		ListAccounts(ctx context.Context, filter model.ListAccountFilter) ([]model.Account, error)
	}

	OutboxMessageDatabasePort interface {
		SaveEvent(ctx context.Context, event model.OutboxEvent) error
		GetEvents(ctx context.Context, limit int64) ([]model.OutboxEvent, error)
		MarkEventAsProcessed(ctx context.Context, eventID uint64) error
	}

	TransactionalDatabasePort interface {
		InsertTransaction(ctx context.Context, tx *model.Transaction) (*model.Transaction, error)
	}

	DatabasePort interface {
		AccountDatabasePort
		OutboxMessageDatabasePort
		TransactionalDatabasePort
	}
)
