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

func (s *TonBeacon) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	log.Debug().Msg("list accounts")

	var filter model.ListAccountFilter
	if len(req.GetWalletIds()) > 0 {
		filter.WalletIDs = lo.ToPtr(req.GetWalletIds())
	}

	if req.IsActive != nil {
		filter.IsClosed = lo.ToPtr(false)
	}

	filter.Offset = int(req.GetOffset())

	if req.GetLimit() == 0 || req.GetLimit() > 1000 {
		req.Limit = 1000
	}
	filter.Limit = int(req.GetLimit())

	accounts, err := s.accountSvc.ListAccounts(ctx, filter)
	if err != nil {
		return listAccountsPbError(codes.Internal, errors.Wrap(err, "list accounts")), nil
	}

	pbAccounts := make([]*pb.Account, 0, len(accounts))
	for _, account := range accounts {
		pbAccount := pb.Account{
			AccountId: account.ID,
			WalletId:  account.WalletID,
			Address:   account.Address.String(),
		}
		pbAccounts = append(pbAccounts, &pbAccount)
	}

	return &pb.ListAccountsResponse{Accounts: pbAccounts}, nil
}

func listAccountsPbError(code codes.Code, err error) *pb.ListAccountsResponse {
	return &pb.ListAccountsResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
