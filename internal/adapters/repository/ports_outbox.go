package repository

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/kriuchkov/tonbeacon/internal/core/model"
)

func (d *DatabaseAdapter) SaveEvent(ctx context.Context, event model.OutboxEvent) error {
	_, err := d.db.NewInsert().Model(fromModelOutboxEvent(event)).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "save event")
	}
	return nil
}

func (d *DatabaseAdapter) GetEvents(ctx context.Context, limit uint64) ([]model.OutboxEvent, error) {
	var events []OutboxEvent
	err := d.db.NewSelect().Model(&events).Limit(int(limit)).For("UPDATE").Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get events")
	}

	var result []model.OutboxEvent
	for i := range events {
		result = append(result, events[i].toModel())
	}

	return result, nil
}

func (d *DatabaseAdapter) MarkEventAsProcessed(ctx context.Context, eventID uint64) error {
	if _, err := d.db.NewUpdate().Model((*OutboxEvent)(nil)).Set("processed = ?", true).
		Where("id = ?", eventID).Exec(ctx); err != nil {
		return errors.Wrap(err, "mark event as processed")
	}
	return nil
}
