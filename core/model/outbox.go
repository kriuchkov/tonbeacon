package model

import (
	"fmt"
	"time"
)

type OutboxEvent struct {
	ID        int64
	EventType EventType
	Payload   []byte // Данные события в JSON
	CreatedAt time.Time
	Processed bool // Флаг обработки
}

func (o OutboxEvent) Key() string {
	return fmt.Sprintf("%s:%d", o.EventType, o.ID)
}
