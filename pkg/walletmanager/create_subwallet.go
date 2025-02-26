package walletmanager

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func (wm *WalletManager) CreateSubWallet(ctx context.Context, subwalletID uint32) (*wallet.Wallet, error) {
	w, err := wallet.FromPrivateKey(wm.api, wm.master.PrivateKey(), wallet.V4R2)
	if err != nil {
		return nil, err
	}

	subwallet, err := w.GetSubwallet(subwalletID)
	if err != nil {
		return nil, errors.Wrap(err, "get subwallet")
	}
	return subwallet, nil
}
