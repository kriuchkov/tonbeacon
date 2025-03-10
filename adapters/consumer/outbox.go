package consumer

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

const (
	defaultInterval = 10 * time.Millisecond
)

type OutboxOptions struct {
	OutboxManager ports.OutboxServicePort       `required:"true"`
	TxManager     ports.DatabaseTransactionPort `required:"true"`
	Writer        *kafka.Writer                 `required:"true"`
	Interval      time.Duration
}

func (o *OutboxOptions) SetDefaults() {
	if o.Interval == 0 {
		o.Interval = defaultInterval
	}
}

type Outbox struct {
	tx        ports.DatabaseTransactionPort
	outboxSvc ports.OutboxServicePort
	writer    *kafka.Writer
	interval  time.Duration
}

func NewOutbox(options OutboxOptions) *Outbox {
	if err := validator.New().Struct(options); err != nil {
		log.Panic().Err(err).Msg("outbox options")
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
	ctx, err := o.tx.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer o.tx.Rollback(ctx) //nolint:errcheck // we don't care about rollback errors

	event, err := o.outboxSvc.GetPendingEvent(ctx)
	if err != nil {
		return errors.Wrap(err, "get pending event")
	}

	if err = o.writer.WriteMessages(ctx, kafka.Message{Key: []byte(event.EventType), Value: event.Payload}); err != nil {
		return errors.Wrap(err, "write message")
	}

	if err = o.outboxSvc.MarkEventAsProcessed(ctx, event.ID); err != nil {
		return errors.Wrap(err, "mark event as processed")
	}

	if err = o.tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}
