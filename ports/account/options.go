package account

import (
	"context"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

type Options struct {
	WalletManager   ports.WalletPort              `required:"true"`
	TxManager       ports.DatabaseTransactionPort `required:"true"`
	DatabaseManager ports.DatabasePort            `required:"true"`
	EventManager    ports.OutboxMessagePort
}

func (o *Options) SetDefaults() {
	if o.EventManager == nil {
		o.EventManager = &emptyOutboxMessagePort{}
	}
}

type emptyOutboxMessagePort struct{}

func (e *emptyOutboxMessagePort) Publish(_ context.Context, _ model.EventType, _ any) error {
	return nil
}
