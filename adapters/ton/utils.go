package ton

import (
	"context"
	"log"
	"math/big"
	"math/rand"

	"github.com/go-faster/errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func TonTransfer(ctx context.Context, from, to *wallet.Wallet, comment string) error {
	var err error

	if from == nil || to == nil || to.Address() == nil {
		return ErrWalletIsEmpty
	}

	var body *cell.Cell
	if comment != "" {
		body, err = buildComment(comment)
		if err != nil {
			return errors.Wrap(err, "build comment")
		}
	}

	return from.Send(ctx, &wallet.Message{
		Mode: 128 + 32, // 128 + 32 send all and destroy
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     to.Address(),
			Amount:      tlb.FromNanoTONU(0),
			Body:        body,
		},
	}, false)
}

func JettonsTransfer(
	ctx context.Context,
	from, to *wallet.Wallet,
	jettonWallet *address.Address,
	forwardAmount tlb.Coins,
	amount Coins,
	comment string,
) error {
	if from == nil || to == nil || to.Address() == nil {
		return errors.New("nil wallet")
	}
	body, err := MakeJettonTransferMessage(to.Address(), to.Address(), amount.BigInt(), forwardAmount, rand.Int63(), comment, "")
	if err != nil {
		return errors.Wrap(err, "make jetton transfer message")
	}

	return from.Send(ctx, &wallet.Message{
		Mode: 128 + 32, // 128 + 32 send all and destroy
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     jettonWallet,
			Amount:      tlb.FromNanoTONU(0),
			Body:        body,
		},
	}, false)
}

func buildComment(comment string) (*cell.Cell, error) {
	root := cell.BeginCell().MustStoreUInt(0, 32)
	if err := root.StoreStringSnake(comment); err != nil {
		return nil, errors.Wrap(err, "store comment")
	}
	return root.EndCell(), nil
}

func MakeJettonTransferMessage(
	destination, responseDest *address.Address,
	amount *big.Int,
	forwardAmount tlb.Coins,
	queryID int64,
	comment string,
	binaryComment string,
) (*cell.Cell, error) {
	var err error

	forwardPayload := cell.BeginCell().EndCell()

	if binaryComment != "" {
		var cell *cell.Cell
		cell, err = decodeBinaryComment(binaryComment)
		if err != nil {
			return nil, errors.Wrap(err, "decode binary comment")
		}

		forwardPayload = cell
	} else if comment != "" {
		forwardPayload, err = buildComment(comment)
		if err != nil {
			return nil, errors.Wrap(err, "build comment")
		}
	}

	payload, err := tlb.ToCell(jetton.TransferPayload{
		QueryID:             uint64(queryID),
		Amount:              tlb.FromNanoTON(amount),
		Destination:         destination,
		ResponseDestination: responseDest,
		CustomPayload:       nil,
		ForwardTONAmount:    forwardAmount,
		ForwardPayload:      forwardPayload,
	})

	if err != nil {
		log.Fatalf("jetton transfer message serialization error: %s", err.Error())
	}

	return payload, nil
}

func decodeBinaryComment(comment string) (*cell.Cell, error) {
	return nil, nil
}
