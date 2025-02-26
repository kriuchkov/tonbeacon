package model

type AccountID = string

type Account struct {
	ID       AccountID
	WalletID uint32
	Address  Address
}

type ListAccountFilter struct {
	WalletIDs *[]uint32
	IsClosed  *bool
	Offset    int
	Limit     int
}
