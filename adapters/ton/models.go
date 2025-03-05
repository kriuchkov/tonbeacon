package ton

import (
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/shopspring/decimal"
)

var _ model.WalletWrapper = &TonWallet{}

type Coins = decimal.Decimal
type APIClientWrapped = ton.APIClientWrapped

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

type shardID struct {
	Workchain int32
	Shard     int64
}

func getShardID(shard *ton.BlockIDExt) shardID {
	return shardID{Workchain: shard.Workchain, Shard: shard.Shard}
}
