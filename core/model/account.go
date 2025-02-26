package model

type AccountID = string

type Account struct {
	ID       AccountID
	WalletID uint32
	Address  Address
}
