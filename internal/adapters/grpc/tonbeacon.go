package grpc

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/kriuchkov/tonbeacon/internal/core/model"
	"github.com/kriuchkov/tonbeacon/internal/core/ports"
	pb "github.com/kriuchkov/tonbeacon/proto/v1"
	"google.golang.org/grpc/codes"
)

type Server struct {
	pb.UnimplementedTonBeaconServer
	accountSvc ports.AccountServicePort
}

func New(account ports.AccountServicePort) *Server {
	return &Server{accountSvc: account}
}

func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	account, err := s.accountSvc.CreateAccount(ctx, req.AccountId)
	if err != nil {
		if errors.Is(err, model.ErrAccountExists) {
			return createAccountPbError(codes.AlreadyExists, err), nil
		}

		return createAccountPbError(codes.Internal, errors.Wrap(err, "create account")), nil
	}

	pbAccount := pb.Account{
		AccountId:   account.ID,
		SubwalletId: account.SubwalletID,
		Address:     account.Address,
	}
	return &pb.CreateAccountResponse{Account: &pbAccount}, nil
}

func createAccountPbError(code codes.Code, err error) *pb.CreateAccountResponse {
	return &pb.CreateAccountResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}

func (s *Server) CloseAccount(ctx context.Context, req *pb.CloseAccountRequest) (*pb.CloseAccountResponse, error) {
	err := s.accountSvc.CloseAccount(ctx, req.AccountId)
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

// GetBalance retrieves the balance for a given account ID.
// It uses the account service to fetch the balance and returns a gRPC response.
//
// Parameters:
//
//	ctx - The context for the request, used for cancellation and deadlines.
//	req - The request containing the account ID for which the balance is to be retrieved.
//
// Returns:
//
//	*pb.GetBalanceResponse - The response containing the balance of the account.
//	error - An error if the balance could not be retrieved, or nil if successful.
//
// Possible errors:
//   - If the account is not found, a gRPC NotFound error is returned.
//   - For other errors, a gRPC Internal error is returned with the wrapped error message.
func (s *Server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	balance, err := s.accountSvc.GetBalance(ctx, req.AccountId)
	if err != nil {
		if errors.Is(err, model.ErrAccountNotFound) {
			return getBalancePbError(codes.NotFound, err), nil
		}
		return getBalancePbError(codes.Internal, errors.Wrap(err, "get balance")), nil
	}

	return &pb.GetBalanceResponse{Balance: int64(balance)}, nil
}

func getBalancePbError(code codes.Code, err error) *pb.GetBalanceResponse {
	return &pb.GetBalanceResponse{Error: &pb.Error{Code: uint32(code), Message: err.Error()}}
}
