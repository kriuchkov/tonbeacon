package model

import "math/big"

type Balance [2]big.Int

type Address string

func (a Address) String() string {
	return string(a)
}

type WalletWrapper interface {
	WalletAddress() Address
}
