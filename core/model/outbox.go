package model

import "time"

type OutboxEvent struct {
	ID        int64
	EventType EventType
	Payload   []byte // Данные события в JSON
	CreatedAt time.Time
	Processed bool // Флаг обработки
}
