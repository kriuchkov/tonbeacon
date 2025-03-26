package model

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type Currency string

func (c Currency) String() string {
	return string(c)
}

const (
	CurrencyTON  Currency = "TON"
	CurrencyUSDT Currency = "USDT"
)

type Amount big.Int

func (a *Amount) String() string {
	if a == nil {
		return "0"
	}

	bigInt := (*big.Int)(a)
	decimalValue := decimal.NewFromBigInt(bigInt, 0)
	result := decimalValue.Div(decimal.NewFromInt(1000000000))
	return result.String()
}

type Balance struct {
	Currency Currency
	Amount   Amount
}

type Address string

func (a Address) String() string {
	return string(a)
}

type WalletWrapper interface {
	WalletAddress() Address
}
