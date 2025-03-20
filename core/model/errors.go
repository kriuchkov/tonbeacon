package model

import "errors"

var (
	ErrAccountExists   = errors.New("account already exists")
	ErrAccountNotFound = errors.New("account not found")
	ErrNoPendingEvents = errors.New("no pending events")
)
