package ton

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

type (
	WalletWrapped interface {
		PrivateKey() ed25519.PrivateKey
		GetSubwallet(subwallet uint32) (*wallet.Wallet, error)
		WalletAddress() *address.Address
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

func (w *WalletAdapter) GetExtraCurrenciesBalance(ctx context.Context, walletID uint32) ([]model.Balance, error) {
	masterBlock, err := w.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get masterchain info")
	}

	subwallet, err := w.masterWallet.GetSubwallet(walletID)
	if err != nil {
		return nil, errors.Wrap(err, "get subwallet")
	}

	adr := subwallet.Address()

	account, err := w.api.WaitForBlock(masterBlock.SeqNo).GetAccount(ctx, masterBlock, adr)
	if err != nil {
		return nil, errors.Wrap(err, "get account")
	}

	log.Debug().Any("account", account).Msg("account")

	if account.IsActive && account.State != nil {
		var currencies []cell.DictKV
		if currencies, err = account.State.ExtraCurrencies.LoadAll(); err != nil {
			return nil, errors.Wrap(err, "load currencies")
		}

		log.Debug().Any("currencies", currencies).Msg("currencies")

		balances := make([]model.Balance, 0, len(currencies))

		for _, kv := range currencies {
			id := kv.Key.MustLoadBigUInt(32)
			amount := kv.Value.MustLoadVarUInt(32)

			log.Debug().Str("id", id.String()).Str("amount", amount.String()).Msg("balance")

			balances = append(balances, model.Balance{*id, *amount})
		}
		return balances, nil
	}
	return nil, nil
}

func (w *WalletAdapter) GetBalance(ctx context.Context, walletID uint32) (uint64, error) {
	subwallet, err := w.masterWallet.GetSubwallet(walletID)
	if err != nil {
		return 0, errors.Wrap(err, "get subwallet")
	}

	masterBlock, err := w.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "get masterchain info")
	}

	balance, err := subwallet.GetBalance(ctx, masterBlock)
	if err != nil {
		return 0, errors.Wrap(err, "get balance")
	}
	return balance.Nano().Uint64(), nil
}

func (w *WalletAdapter) SendWaitTransaction(ctx context.Context, walletID uint32) error {
	subwallet, err := w.masterWallet.GetSubwallet(walletID)
	if err != nil {
		return errors.Wrap(err, "get subwallet")
	}

	masterBlock, err := w.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "get masterchain info")
	}

	balance, err := subwallet.GetBalance(ctx, masterBlock)
	if err != nil {
		return errors.Wrap(err, "get balance")
	}

	fmt.Print("balance: ", balance)
	return errors.New("account is not active")
}

func (w *WalletAdapter) TransferToMainWallet(ctx context.Context, walletID uint32, amount uint64) error {
	subwallet, err := w.masterWallet.GetSubwallet(walletID)
	if err != nil {
		return errors.Wrap(err, "get subwallet")
	}

	// Get the current block info
	block, err := w.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "get masterchain info")
	}

	// Check if the subwallet has enough balance
	balance, err := subwallet.GetBalance(ctx, block)
	if err != nil {
		return errors.Wrap(err, "get balance")
	}

	if balance.Nano().Uint64() <= amount {
		return errors.New("insufficient balance in subwallet")
	}

	coinAmount := tlb.MustFromNano(big.NewInt(int64(amount)), 0)
	if err = subwallet.TransferNoBounce(ctx, w.masterWallet.WalletAddress(), coinAmount, ""); err != nil {
		return errors.Wrap(err, "transfer funds to main wallet")
	}

	log.Info().
		Uint32("from_wallet_id", walletID).
		Uint64("amount", amount).
		Str("to_address", w.masterWallet.WalletAddress().String()).
		Msg("transferring funds to main wallet")

	return nil
}
