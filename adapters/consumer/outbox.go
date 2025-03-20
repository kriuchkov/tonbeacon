package consumer

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

const (
	defaultInterval = 10 * time.Millisecond
)

type OutboxWriter interface {
	SendMessage(key string, value []byte) (partition int32, offset int64, err error)
}

type OutboxOptions struct {
	OutboxManager ports.OutboxServicePort             `validate:"required"`
	TxManager     ports.DatabaseWithinTransactionPort `validate:"required"`
	Writer        OutboxWriter                        `validate:"required"`
	Interval      time.Duration
}

func (o *OutboxOptions) SetDefaults() {
	if o.Interval == 0 {
		o.Interval = defaultInterval
	}
}

type Outbox struct {
	tx        ports.DatabaseWithinTransactionPort
	outboxSvc ports.OutboxServicePort
	writer    OutboxWriter
	interval  time.Duration
}

func NewOutbox(options OutboxOptions) *Outbox {
	options.SetDefaults()

	if err := validator.New().Struct(options); err != nil {
		panic(err.Error())
	}

	return &Outbox{
		tx:        options.TxManager,
		outboxSvc: options.OutboxManager,
		writer:    options.Writer,
		interval:  options.Interval,
	}
}

func (o *Outbox) Consumer(ctx context.Context) {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := o.process(context.Background()); err != nil {
				log.Warn().Err(err).Msg("process outbox")
			}

			<-ticker.C
		}
	}
}

func (o *Outbox) process(ctx context.Context) error {
	err := o.tx.WithInTransaction(ctx, func(ctx context.Context) error {
		event, err := o.outboxSvc.GetPendingEvent(ctx)
		if err != nil {
			if errors.Is(err, model.ErrNoPendingEvents) {
				return nil
			}

			return errors.Wrap(err, "get pending event")
		}

		if _, _, err = o.writer.SendMessage(event.Key(), event.Payload); err != nil {
			return errors.Wrap(err, "send message")
		}

		if err = o.outboxSvc.MarkEventAsProcessed(ctx, event.ID); err != nil {
			return errors.Wrap(err, "mark event as processed")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "process transaction")
	}
	return nil
}
