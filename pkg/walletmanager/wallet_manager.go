package walletmanager

import (
	"context"

	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

// Интерфейс для наблюдателя транзакций.
type Repository interface {
	IsSubWalletExist(ctx context.Context, walletID uint32) (bool, error)
}

type WalletManager struct {
	api        *ton.APIClient
	master     *wallet.Wallet
	repository Repository
}

func New(masterSeed []string, tonAddr, tonKey string) (*WalletManager, error) {
	return &WalletManager{}, nil
}
