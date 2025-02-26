package ports

import (
	"context"

	"github.com/kriuchkov/tonbeacon/core/model"
)

type WalletPort interface {
	CreateWallet(ctx context.Context, walletID uint32) (model.WalletWrapper, error)
}

type AccountDatabasePort interface {
	IsAccountExists(ctx context.Context, accountID string) (bool, error)
	InsertAccount(ctx context.Context, accountID string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) error
	CloseAccount(ctx context.Context, accountID string) error
}

type OutboxMessageDatabasePort interface {
	SaveEvent(ctx context.Context, event model.OutboxEvent) error
	GetEvents(ctx context.Context, limit uint64) ([]model.OutboxEvent, error)
	MarkEventAsProcessed(ctx context.Context, eventID uint64) error
}

type DatabaseTransactionPort interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DatabasePort interface {
	AccountDatabasePort
	OutboxMessageDatabasePort
}

type OutboxMessagePort interface {
	Publish(ctx context.Context, eventType model.EventType, payload any) error
}
