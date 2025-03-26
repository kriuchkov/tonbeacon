package model

type AccountID = string

type Account struct {
	ID       AccountID
	WalletID uint32
	Address  Address
}

type ListAccountFilter struct {
	ID        *AccountID
	WalletIDs *[]uint32
	IsClosed  *bool
	Address   *string
	Offset    int
	Limit     int
}
