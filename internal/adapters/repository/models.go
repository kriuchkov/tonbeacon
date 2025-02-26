package repository

import (
	"time"

	"github.com/kriuchkov/tonbeacon/internal/core/model"
	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID          string `bun:"id,pk"`
	SubwalletID uint32 `bun:"subwallet_id,unique,notnull"`
	Address     string `bun:"address"`
	IsClosed    bool   `bun:"is_closed"`
}

func (a *Account) toModel() *model.Account {
	return &model.Account{
		ID:          a.ID,
		SubwalletID: a.SubwalletID,
		Address:     a.Address,
	}
}

func fromModelAccount(account *model.Account) *Account {
	return &Account{
		ID:          account.ID,
		SubwalletID: account.SubwalletID,
		Address:     account.Address,
	}
}

type OutboxEvent struct {
	bun.BaseModel `bun:"table:outbox_events"`

	ID        uint64    `bun:"id,pk"`
	EventType string    `bun:"event_type"`
	Payload   string    `bun:"payload"`
	CreatedAt time.Time `bun:"created_at"`
	Processed bool      `bun:"processed"`
}

func (e *OutboxEvent) toModel() model.OutboxEvent {
	return model.OutboxEvent{
		ID:        int64(e.ID),
		EventType: model.EventType(e.EventType),
		Payload:   []byte(e.Payload),
		CreatedAt: e.CreatedAt,
		Processed: e.Processed,
	}
}

func fromModelOutboxEvent(event model.OutboxEvent) *OutboxEvent {
	return &OutboxEvent{
		ID:        uint64(event.ID),
		EventType: string(event.EventType),
		Payload:   string(event.Payload),
		CreatedAt: event.CreatedAt,
		Processed: event.Processed,
	}
}
