package ton

import (
	"context"
	"crypto/ed25519"

	"github.com/go-faster/errors"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

type (
	APIClientWrapped = ton.APIClientWrapped

	WalletWrapped interface {
		PrivateKey() ed25519.PrivateKey
		GetSubwallet(subwallet uint32) (*wallet.Wallet, error)
	}
)

var _ ports.WalletPort = (*WalletAdapter)(nil)
var _ WalletWrapped = &wallet.Wallet{}

type WalletAdapter struct {
	api          APIClientWrapped
	masterWallet WalletWrapped
}

func NewWalletAdapter(api APIClientWrapped, masterWallet WalletWrapped) *WalletAdapter {
	return &WalletAdapter{api: api, masterWallet: masterWallet}
}

func (w *WalletAdapter) CreateWallet(ctx context.Context, walletID uint32) (model.WalletWrapper, error) {
	subwallet, err := w.masterWallet.GetSubwallet(walletID)
	if err != nil {
		return nil, errors.Wrap(err, "get subwallet")
	}
	return fromUtilsWallet(subwallet), nil
}
