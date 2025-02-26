package ton

import (
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/kriuchkov/tonbeacon/core/model"
)

var _ model.WalletWrapper = &TonWallet{}

type TonWallet struct {
	Address model.Address
}

func (w *TonWallet) WalletAddress() model.Address {
	return w.Address
}

func fromUtilsWallet(w *wallet.Wallet) *TonWallet {
	return &TonWallet{
		Address: model.Address(w.WalletAddress().String()),
	}
}
