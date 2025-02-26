package repository

import (
	"time"

	"github.com/uptrace/bun"

	"github.com/kriuchkov/tonbeacon/core/model"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID         string  `bun:"id,pk"`
	WalletID   uint32  `bun:"wallet_id,unique,notnull"`
	TonAddress *string `bun:"ton_address"`
	IsClosed   bool    `bun:"is_closed"`
}

func (a *Account) toModel() *model.Account {
	var address model.Address
	if a.TonAddress != nil {
		address = model.Address(*a.TonAddress)
	}
	return &model.Account{
		ID:       a.ID,
		WalletID: a.WalletID,
		Address:  address,
	}
}

func fromModelAccount(account *model.Account) *Account {
	var tonAddress *string
	if account.Address != "" {
		addrStr := account.Address.String()
		tonAddress = &addrStr
	}
	return &Account{
		ID:         account.ID,
		WalletID:   account.WalletID,
		TonAddress: tonAddress,
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
