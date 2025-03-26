package grpc

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"

	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
)

func (s *TonBeacon) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	log.Debug().Str("account_id", req.GetAccountId()).Msg("get account")

	var listReq model.ListAccountFilter
	if req.AccountId != nil {
		listReq.ID = lo.ToPtr(req.GetAccountId())
	}

	if req.WalletId != nil {
		var walletID []uint32
		walletID = append(walletID, req.GetWalletId())
		listReq.WalletIDs = &walletID
	}

	if req.Address != nil {
		listReq.Address = lo.ToPtr(req.GetAddress())
	}

	accounts, err := s.accountSvc.ListAccounts(ctx, listReq)
	if err != nil {
		if errors.Is(err, model.ErrAccountNotFound) {
			return getAccountPbError(codes.NotFound, err), nil
		}

		log.Err(err).Msg("get account")
		return getAccountPbError(codes.Internal, errors.Wrap(err, "get account")), nil
	}

	if len(accounts) == 0 {
		return getAccountPbError(codes.NotFound, model.ErrAccountNotFound), nil
	}

	account := accounts[0]

	pbAccount := pb.Account{
		AccountId: account.ID,
		WalletId:  account.WalletID,
		Address:   account.Address.String(),
	}
	return &pb.GetAccountResponse{Account: &pbAccount}, nil
}

func getAccountPbError(code codes.Code, err error) *pb.GetAccountResponse {
	return &pb.GetAccountResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
