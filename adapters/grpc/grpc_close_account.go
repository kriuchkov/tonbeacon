package grpc

import (
	"context"

	"github.com/go-faster/errors"
	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/model"
	"google.golang.org/grpc/codes"
)

func (s *TonBeacon) CloseAccount(ctx context.Context, req *pb.CloseAccountRequest) (*pb.CloseAccountResponse, error) {
	err := s.accountSvc.CloseAccount(ctx, req.GetAccountId())
	if err != nil {
		if errors.Is(err, model.ErrAccountNotFound) {
			return closeAccountPbError(codes.NotFound, err), nil
		}
		return closeAccountPbError(codes.Internal, errors.Wrap(err, "close account")), nil
	}
	return &pb.CloseAccountResponse{}, nil
}

func closeAccountPbError(code codes.Code, err error) *pb.CloseAccountResponse {
	return &pb.CloseAccountResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
