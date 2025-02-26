package model

type AccountID = string

type Account struct {
	ID          AccountID
	SubwalletID uint32
	Address     string
}
