package grpc

import (
	"context"
	"net"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

type TonBeacon struct {
	pb.UnimplementedTonBeaconServer
	accountSvc ports.AccountServicePort
	server     *grpc.Server
}

func NewTonBeacon(account ports.AccountServicePort) *TonBeacon {
	return &TonBeacon{accountSvc: account, server: grpc.NewServer()}
}

func (s *TonBeacon) Run(lis net.Listener) error {
	pb.RegisterTonBeaconServer(s.server, s)
	reflection.Register(s.server)
	return s.server.Serve(lis)
}

func (s *TonBeacon) Stop() {
	s.server.Stop()
}

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
