package grpc

import (
	"context"

	"github.com/go-faster/errors"
	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

func (s *TonBeacon) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	log.Debug().Str("account_id", req.GetAccountId()).Msg("get balance")

	balance, err := s.accountSvc.GetBalance(ctx, req.GetAccountId())
	if err != nil {
		if errors.Is(err, model.ErrAccountNotFound) {
			return getBalancePbError(codes.NotFound, err), nil
		}
		return getBalancePbError(codes.Internal, errors.Wrap(err, "get balance")), nil
	}

	var balances []*pb.Tokens
	for _, b := range balance {
		balances = append(balances, &pb.Tokens{Symbol: b.Currency.String(), Amount: b.Amount.String()})
	}
	return &pb.GetBalanceResponse{Tokens: balances}, nil
}

func getBalancePbError(code codes.Code, err error) *pb.GetBalanceResponse {
	return &pb.GetBalanceResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
