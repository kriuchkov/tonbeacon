package account_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kriuchkov/tonbeacon/core/model"
	portsmocks "github.com/kriuchkov/tonbeacon/core/ports/mocks"
	"github.com/kriuchkov/tonbeacon/ports/account"
)

var _ model.WalletWrapper = (*TestWalletWrapper)(nil)

type TestWalletWrapper struct {
	address model.Address
}

func (t *TestWalletWrapper) WalletAddress() model.Address {
	return t.address
}

func TestAccount_CreateAccount(t *testing.T) {
	type isAccountExistsMock struct {
		callCount     int
		accountID     string
		exists        bool
		expectedError error
	}

	type beginMock struct {
		callCount     int
		expectedError error
	}

	type rollbackMock struct {
		callCount     int
		expectedError error
	}

	type commitMock struct {
		callCount     int
		expectedError error
	}

	type insertAccountMock struct {
		callCount     int
		accountID     string
		result        *model.Account
		expectedError error
	}

	type createWalletMock struct {
		callCount     int
		walletID      uint32
		result        model.WalletWrapper
		expectedError error
	}

	type updateAccountMock struct {
		callCount     int
		account       *model.Account
		expectedError error
	}

	type publishEventMock struct {
		callCount     int
		eventType     model.EventType
		payload       any
		expectedError error
	}

	tests := []struct {
		name           string
		accountID      string
		expectedError  error
		expectedResult *model.Account

		isAccountExistsMock
		beginMock
		rollbackMock
		commitMock
		insertAccountMock
		createWalletMock
		updateAccountMock
		publishEventMock
	}{
		{
			name:         "successful account creation",
			accountID:    "test-account",
			beginMock:    beginMock{callCount: 1},
			rollbackMock: rollbackMock{callCount: 1},
			commitMock:   commitMock{callCount: 1},
			isAccountExistsMock: isAccountExistsMock{
				callCount: 1,
				accountID: "test-account",
			},
			insertAccountMock: insertAccountMock{
				callCount: 1,
				accountID: "test-account",
				result:    &model.Account{ID: "test-account", WalletID: 10},
			},
			createWalletMock: createWalletMock{
				callCount: 1,
				walletID:  10,
				result:    &TestWalletWrapper{address: model.Address("test-address")},
			},
			updateAccountMock: updateAccountMock{
				callCount: 1,
				account:   &model.Account{ID: "test-account", WalletID: 10, Address: "test-address"},
			},
			publishEventMock: publishEventMock{
				callCount: 1,
				eventType: model.AccountCreated,
				payload:   &model.Account{ID: "test-account", WalletID: 10, Address: "test-address"},
			},
			expectedResult: &model.Account{ID: "test-account", WalletID: 10, Address: "test-address"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Transaction manager mock
			mockTx := portsmocks.NewMockDatabaseTransactionPort(t)
			if tt.beginMock.callCount > 0 {
				param := tt.beginMock
				mockTx.On("Begin", ctx).Return(context.Background(), param.expectedError).Times(param.callCount)
			}
			if tt.rollbackMock.callCount > 0 {
				param := tt.rollbackMock
				mockTx.On("Rollback", ctx).Return(param.expectedError).Times(param.callCount)
			}
			if tt.commitMock.callCount > 0 {
				param := tt.commitMock
				mockTx.On("Commit", ctx).Return(param.expectedError).Times(param.callCount)
			}

			// Database manager mock
			mockDB := portsmocks.NewMockDatabasePort(t)
			if tt.isAccountExistsMock.callCount > 0 {
				param := tt.isAccountExistsMock
				mockDB.On("IsAccountExists", ctx, param.accountID).Return(param.exists, param.expectedError).Times(param.callCount)
			}
			if tt.insertAccountMock.callCount > 0 {
				param := tt.insertAccountMock
				mockDB.On("InsertAccount", ctx, param.accountID).Return(param.result, param.expectedError).Times(param.callCount)
			}
			if tt.updateAccountMock.callCount > 0 {
				param := tt.updateAccountMock
				mockDB.On("UpdateAccount", ctx, param.account).Return(param.expectedError).Times(param.callCount)
			}

			// Wallet manager mock
			mockWM := portsmocks.NewMockWalletPort(t)
			if tt.createWalletMock.callCount > 0 {
				param := tt.createWalletMock
				mockWM.On("CreateWallet", ctx, param.walletID).Return(param.result, param.expectedError).Times(param.callCount)
			}

			// Event manager mock
			mockEM := portsmocks.NewMockOutboxMessagePort(t)
			if tt.publishEventMock.callCount > 0 {
				param := tt.publishEventMock
				mockEM.On("Publish", ctx, param.eventType, param.payload).Return(param.expectedError).Times(param.callCount)
			}

			accountService := account.New(account.Options{
				WalletManager:   mockWM,
				TxManager:       mockTx,
				DatabaseManager: mockDB,
				EventManager:    mockEM,
			})

			result, err := accountService.CreateAccount(context.Background(), tt.accountID)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedResult, result)

			mockDB.AssertExpectations(t)
			mockWM.AssertExpectations(t)
			mockTx.AssertExpectations(t)
			mockEM.AssertExpectations(t)
		})
	}
}
