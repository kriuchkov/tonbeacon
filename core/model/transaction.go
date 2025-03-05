package model

import "time"

type Transaction struct {
	ID             string
	Sender         string
	Receiver       string
	Amount         float64
	BlockID        string
	CreatedAt      time.Time
	SenderIsOurs   bool
	ReceiverIsOurs bool
}

func UnmarshalTransaction(data []byte) (Transaction, error) {
	return Transaction{}, nil
}
