package grpc

import (
	"context"

	"github.com/go-faster/errors"
	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

func (s *TonBeacon) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	log.Debug().Str("account_id", req.GetAccountId()).Msg("create account")

	account, err := s.accountSvc.CreateAccount(ctx, req.GetAccountId())
	if err != nil {
		if errors.Is(err, model.ErrAccountExists) {
			return createAccountPbError(codes.AlreadyExists, err), nil
		}

		log.Err(err).Msg("create account")
		return createAccountPbError(codes.Internal, errors.Wrap(err, "create account")), nil
	}

	pbAccount := pb.Account{
		AccountId: account.ID,
		WalletId:  account.WalletID,
		Address:   account.Address.String(),
	}
	return &pb.CreateAccountResponse{Account: &pbAccount}, nil
}

func createAccountPbError(code codes.Code, err error) *pb.CreateAccountResponse {
	return &pb.CreateAccountResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
