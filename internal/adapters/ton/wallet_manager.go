package ton

import (
	"context"

	"github.com/kriuchkov/tonbeacon/internal/core/ports"
	"github.com/kriuchkov/tonbeacon/pkg/walletmanager"
)

var _ ports.WalletPort = (*WalletManager)(nil)

type WalletManager struct {
	manager *walletmanager.WalletManager
}

func NewWalletManager() *WalletManager {
	return &WalletManager{}
}

func (w *WalletManager) CreateWallet(ctx context.Context, walletID uint32) (string, error) {
	walletAddress, err := w.manager.CreateSubWallet(ctx, walletID)
	if err != nil {
		return "", err
	}
	return walletAddress.Address().String(), nil
}
