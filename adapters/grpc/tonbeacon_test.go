package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	grpcAdapter "github.com/kriuchkov/tonbeacon/adapters/grpc"
	pb "github.com/kriuchkov/tonbeacon/api/grpc/v1"
	"github.com/kriuchkov/tonbeacon/core/model"
	portsmocks "github.com/kriuchkov/tonbeacon/core/ports/mocks"
)

func TestTonBeacon_CreateAccount(t *testing.T) {
	type mockAccountServiceCall struct {
		calls           int
		accountID       string
		expectedAccount *model.Account
		expectedError   error
	}
	tests := []struct {
		name             string
		req              *pb.CreateAccountRequest
		expectedResponse *pb.CreateAccountResponse
		expectedError    error

		mockAccountServiceCall
	}{
		{
			name: "successful account creation",
			req:  &pb.CreateAccountRequest{AccountId: "test-account-123"},
			mockAccountServiceCall: mockAccountServiceCall{
				calls:     1,
				accountID: "test-account-123",
				expectedAccount: &model.Account{
					ID:       "test-account-123",
					WalletID: 42,
					Address:  "EQD-oGD82Hyr7P6GWNCyMNRj_qA_7dF93NUzeFhVvrxtD510",
				},
			},
			expectedResponse: &pb.CreateAccountResponse{
				Account: &pb.Account{
					AccountId: "test-account-123",
					WalletId:  42,
					Address:   "EQD-oGD82Hyr7P6GWNCyMNRj_qA_7dF93NUzeFhVvrxtD510",
				},
			},
		},
		{
			name: "account already exists",
			req:  &pb.CreateAccountRequest{AccountId: "existing-account"},
			mockAccountServiceCall: mockAccountServiceCall{
				calls:         1,
				accountID:     "existing-account",
				expectedError: model.ErrAccountExists,
			},
			expectedResponse: &pb.CreateAccountResponse{
				Error: &pb.Error{
					Code:    uint32(codes.AlreadyExists),
					Message: model.ErrAccountExists.Error(),
				},
			},
		},
		{
			name: "internal error during account creation",
			req:  &pb.CreateAccountRequest{AccountId: "error-account"},
			mockAccountServiceCall: mockAccountServiceCall{
				calls:         1,
				accountID:     "error-account",
				expectedError: errors.New("database connection failed"),
			},
			expectedResponse: &pb.CreateAccountResponse{
				Error: &pb.Error{
					Code:    uint32(codes.Internal),
					Message: "create account: database connection failed",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccountSvc := portsmocks.NewMockAccountServicePort(t)
			if tt.mockAccountServiceCall.calls > 0 {
				param := tt.mockAccountServiceCall
				mockAccountSvc.On("CreateAccount", mock.Anything, param.accountID).
					Return(param.expectedAccount, param.expectedError).
					Times(param.calls)
			}

			tonBeacon := grpcAdapter.NewTonBeacon(mockAccountSvc)

			resp, err := tonBeacon.CreateAccount(context.Background(), tt.req)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedResponse, resp)

			mockAccountSvc.AssertExpectations(t)
		})
	}
}
