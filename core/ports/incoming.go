package ports

import (
	"context"

	"github.com/kriuchkov/tonbeacon/core/model"
)

type AccountServicePort interface {
	CreateAccount(ctx context.Context, accountID model.AccountID) (*model.Account, error)
	GetBalance(ctx context.Context, accountID model.AccountID) (uint64, error)
	CloseAccount(ctx context.Context, accountID model.AccountID) error
}

type OutboxServicePort interface {
	GetPendingEvent(ctx context.Context) (*model.OutboxEvent, error)
	MarkEventAsProcessed(ctx context.Context, eventID int64) error
}
