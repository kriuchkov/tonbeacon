package repository

import (
	"time"

	"github.com/uptrace/bun"

	"github.com/kriuchkov/tonbeacon/core/model"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            string  `bun:"id,pk"`
	WalletID      uint32  `bun:"wallet_id,unique,notnull,default"`
	TonAddress    *string `bun:"ton_address"`
	IsClosed      bool    `bun:"is_closed"`
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
		TonAddress: tonAddress,
		WalletID:   account.WalletID,
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

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	ID             string    `bun:"id,pk"`
	Sender         string    `bun:"sender"`
	Receiver       string    `bun:"receiver"`
	Amount         float64   `bun:"amount"`
	BlockID        string    `bun:"block_id"`
	CreatedAt      time.Time `bun:"created_at"`
	SenderIsOurs   bool      `bun:"sender_is_ours"`
	ReceiverIsOurs bool      `bun:"receiver_is_ours"`
}

func (t *Transaction) toModel() *model.Transaction {
	return &model.Transaction{
		ID:             t.ID,
		Sender:         t.Sender,
		Receiver:       t.Receiver,
		Amount:         t.Amount,
		BlockID:        t.BlockID,
		CreatedAt:      t.CreatedAt,
		SenderIsOurs:   t.SenderIsOurs,
		ReceiverIsOurs: t.ReceiverIsOurs,
	}
}

func fromModelTransaction(transaction *model.Transaction) *Transaction {
	return &Transaction{
		ID:             transaction.ID,
		Sender:         transaction.Sender,
		Receiver:       transaction.Receiver,
		Amount:         transaction.Amount,
		BlockID:        transaction.BlockID,
		CreatedAt:      transaction.CreatedAt,
		SenderIsOurs:   transaction.SenderIsOurs,
		ReceiverIsOurs: transaction.ReceiverIsOurs,
	}
}
