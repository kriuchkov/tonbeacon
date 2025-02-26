package repository

import (
	"github.com/kriuchkov/tonbeacon/internal/core/ports"
	"github.com/uptrace/bun"
)

var _ ports.DatabasePort = (*DatabaseAdapter)(nil)

type DatabaseAdapter struct {
	db *bun.DB
}

func New(db *bun.DB) (*DatabaseAdapter, error) {
	return &DatabaseAdapter{db: db}, nil
}
