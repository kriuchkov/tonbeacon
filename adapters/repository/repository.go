package repository

import (
	"github.com/uptrace/bun"

	"github.com/kriuchkov/tonbeacon/core/ports"
)

var _ ports.DatabasePort = (*DatabaseAdapter)(nil)

type DatabaseAdapter struct {
	TxRepository
}

func New(db *bun.DB) *DatabaseAdapter {
	return &DatabaseAdapter{TxRepository{dbConn: db}}
}
