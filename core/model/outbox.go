package model

import (
	"fmt"
	"time"
)

type OutboxEvent struct {
	ID        int64
	EventType EventType
	Payload   []byte
	CreatedAt time.Time
	Processed bool
}

func (o OutboxEvent) Key() string {
	return fmt.Sprintf("%s:%d", o.EventType, o.ID)
}
