package transaction

import (
	"context"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

const (
	// defaultUpdateInterval is the default interval for updating accounts.
	defaultUpdateInterval = 5 * time.Second
)

type Options struct {
	AccountPort     ports.AccountDatabasePort       `required:"true"`
	TxPort          ports.DatabaseTransactionPort   `required:"true"`
	TransactionPort ports.TransactionalDatabasePort `required:"true"`
	Interval        time.Duration
}

func (o *Options) SetDefaults() {
	if o.Interval == 0 {
		o.Interval = defaultUpdateInterval
	}
}

type Transaction struct {
	mx          sync.RWMutex
	accountList map[model.Address]*model.Account
	interval    time.Duration

	// ports
	txPort      ports.DatabaseTransactionPort
	accountPort ports.AccountDatabasePort
	transaction ports.TransactionalDatabasePort
}

func New(ctx context.Context, opts *Options) *Transaction {
	if err := validator.New().Struct(opts); err != nil {
		log.Panic().Err(err).Msg("transaction invalid options")
	}

	t := &Transaction{
		accountPort: opts.AccountPort,
		txPort:      opts.TxPort,
		transaction: opts.TransactionPort,
		accountList: make(map[model.Address]*model.Account),
		interval:    opts.Interval,
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		if err := t.update(ctx, done); err != nil {
			log.Error().Err(err).Msg("update")
		}
	}()

	<-done
	return t
}

func (t *Transaction) update(ctx context.Context, ch chan struct{}) error {
	if err := t.updateAccounts(ctx); err != nil {
		return errors.Wrap(err, "update accounts")
	}

	ch <- struct{}{}

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := t.updateAccounts(ctx); err != nil {
				log.Error().Err(err).Msg("update accounts")
			}
		}
	}
}

func (t *Transaction) updateAccounts(ctx context.Context) error {
	accountList, err := t.accountPort.ListAccounts(ctx, model.ListAccountFilter{IsClosed: lo.ToPtr(false)})
	if err != nil {
		return errors.Wrap(err, "list accounts")
	}

	t.mx.Lock()
	defer t.mx.Unlock()
	for _, account := range accountList {
		t.accountList[account.Address] = &account
	}
	return nil
}

func (t *Transaction) Handle(ctx context.Context, message []byte) error {
	tx, err := model.UnmarshalTransaction(message)
	if err != nil {
		return errors.Wrap(err, "unmarshal tx")
	}

	t.mx.RLock()
	defer t.mx.RUnlock()

	accounts := []*model.Account{
		t.accountList[model.Address(tx.Sender)],
		t.accountList[model.Address(tx.Receiver)],
	}

	if ok := lo.ContainsBy(accounts, func(i *model.Account) bool { return i != nil }); !ok {
		return model.ErrAccountNotFound
	}

	log.Debug().
		Str("sender", tx.Sender).Str("receiver", tx.Receiver).Str("txid", tx.ID).Float64("amount", tx.Amount).
		Msg("processing relevant transaction")

	err = t.txPort.WithInTransaction(ctx, func(ctx context.Context) error {
		if _, err = t.transaction.InsertTransaction(ctx, &tx); err != nil {
			return errors.Wrap(err, "save tx")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "handle tx")
	}
	return nil
}
