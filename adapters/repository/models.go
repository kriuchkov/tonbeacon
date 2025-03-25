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

	ID        int64     `bun:"id,pk"`
	EventType string    `bun:"event_type"`
	Payload   string    `bun:"payload"`
	CreatedAt time.Time `bun:"created_at"`
	Processed bool      `bun:"processed"`
}

func (e *OutboxEvent) toModel() model.OutboxEvent {
	return model.OutboxEvent{
		ID:        e.ID,
		EventType: model.EventType(e.EventType),
		Payload:   []byte(e.Payload),
		CreatedAt: e.CreatedAt,
		Processed: e.Processed,
	}
}

func fromModelOutboxEvent(event model.OutboxEvent) *OutboxEvent {
	return &OutboxEvent{
		ID:        event.ID,
		EventType: string(event.EventType),
		Payload:   string(event.Payload),
		CreatedAt: event.CreatedAt,
		Processed: event.Processed,
	}
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	// Transaction identifiers
	ID          int64  `bun:"id,pk,autoincrement"`
	AccountAddr string `bun:"account_addr"`
	LT          int64  `bun:"lt"`
	PrevTxHash  string `bun:"prev_tx_hash"`
	PrevTxLT    int64  `bun:"prev_tx_lt"`

	// Address information
	Sender   string `bun:"sender"`
	Receiver string `bun:"receiver"`

	// Financial information
	Amount    float64 `bun:"amount"`
	TotalFees float64 `bun:"total_fees"`
	ExitCode  int     `bun:"exit_code"`
	Success   bool    `bun:"success"`

	// Message information
	MessageType string `bun:"message_type"`
	Bounce      bool   `bun:"bounce"`
	Bounced     bool   `bun:"bounced"`
	Body        string `bun:"body"`

	// State information
	BlockID       string    `bun:"block_id"`
	CreatedAt     time.Time `bun:"created_at"`
	AccountStatus string    `bun:"account_status"`

	// Extra info
	ComputeGasUsed int    `bun:"compute_gas_used"`
	Description    string `bun:"description"`
}

func (t *Transaction) toModel() *model.Transaction {
	return &model.Transaction{
		AccountAddr:    t.AccountAddr,
		LT:             t.LT,
		PrevTxHash:     t.PrevTxHash,
		PrevTxLT:       t.PrevTxLT,
		Sender:         t.Sender,
		Receiver:       t.Receiver,
		Amount:         t.Amount,
		TotalFees:      t.TotalFees,
		ExitCode:       t.ExitCode,
		Success:        t.Success,
		MessageType:    t.MessageType,
		Bounce:         t.Bounce,
		Bounced:        t.Bounced,
		Body:           t.Body,
		BlockID:        t.BlockID,
		CreatedAt:      t.CreatedAt,
		AccountStatus:  t.AccountStatus,
		ComputeGasUsed: t.ComputeGasUsed,
		Description:    t.Description,
	}
}

func fromModelTransaction(transaction *model.Transaction) *Transaction {
	return &Transaction{
		AccountAddr:    transaction.AccountAddr,
		LT:             transaction.LT,
		PrevTxHash:     transaction.PrevTxHash,
		PrevTxLT:       transaction.PrevTxLT,
		Sender:         transaction.Sender,
		Receiver:       transaction.Receiver,
		Amount:         transaction.Amount,
		TotalFees:      transaction.TotalFees,
		ExitCode:       transaction.ExitCode,
		Success:        transaction.Success,
		MessageType:    transaction.MessageType,
		Bounce:         transaction.Bounce,
		Bounced:        transaction.Bounced,
		Body:           transaction.Body,
		BlockID:        transaction.BlockID,
		CreatedAt:      transaction.CreatedAt,
		AccountStatus:  transaction.AccountStatus,
		ComputeGasUsed: transaction.ComputeGasUsed,
		Description:    transaction.Description,
	}
}
