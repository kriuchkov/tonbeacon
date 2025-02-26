package model

type EventType string

const (
	AccountCreated EventType = "account_created"
	AccountClosed  EventType = "account_closed"
)
