package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-faster/errors"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

// Outbox is a service that allows to store events that should be processed by external services.
// It is used to implement the outbox pattern.
var _ ports.OutboxServicePort = (*Outbox)(nil)

type Outbox struct {
	database ports.OutboxMessageDatabasePort
}

func New(db ports.OutboxMessageDatabasePort) *Outbox {
	return &Outbox{database: db}
}

func (s *Outbox) AddEvent(ctx context.Context, eventType model.EventType, payload any) error {
	var payloadBytes []byte
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return err
	}

	event := model.OutboxEvent{
		EventType: eventType,
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
		Processed: false,
	}
	return s.database.SaveEvent(ctx, event)
}

func (s *Outbox) GetPendingEvent(ctx context.Context) (*model.OutboxEvent, error) {
	events, err := s.database.GetEvents(ctx, 1)
	if err != nil {
		return nil, errors.Wrap(err, "get events")
	}

	if len(events) == 0 {
		return nil, nil
	}
	return &events[0], nil
}

func (s *Outbox) MarkEventAsProcessed(ctx context.Context, eventID int64) error {
	return s.database.MarkEventAsProcessed(ctx, uint64(eventID))
}
