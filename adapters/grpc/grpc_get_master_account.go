package grpc

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"

	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *TonBeacon) GetMasterAccount(ctx context.Context, _ *emptypb.Empty) (*pb.GetAccountResponse, error) {
	log.Debug().Msg("get master account")

	account, err := s.accountSvc.MasterAccount(ctx)
	if err != nil {
		log.Err(err).Msg("get master account")
		return getMasterAccountPbError(codes.Internal, errors.Wrap(err, "get master account")), nil
	}

	pbAccount := pb.Account{
		AccountId: account.ID,
		WalletId:  account.WalletID,
		Address:   account.Address.String(),
	}
	return &pb.GetAccountResponse{Account: &pbAccount}, nil
}

func getMasterAccountPbError(code codes.Code, err error) *pb.GetAccountResponse {
	return &pb.GetAccountResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
