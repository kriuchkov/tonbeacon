package model

type Address string

func (a Address) String() string {
	return string(a)
}

type WalletWrapper interface {
	WalletAddress() Address
}
