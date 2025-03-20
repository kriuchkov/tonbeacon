package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kriuchkov/tonbeacon/core/model"
	portsmocks "github.com/kriuchkov/tonbeacon/core/ports/mocks"
)

func TestTransaction_Update(t *testing.T) {
	t.Parallel()

	type mockListAccountsCall struct {
		calls       int
		filter      model.ListAccountFilter
		accountList []model.Account
		expectError error
	}

	tests := []struct {
		name        string
		accountList map[model.Address]*model.Account

		mockListAccountsCall
	}{
		{
			name: "successful",
			accountList: map[model.Address]*model.Account{
				"addr1": {ID: "1", WalletID: 1, Address: "addr1"},
				"addr2": {ID: "2", WalletID: 2, Address: "addr2"},
			},
			mockListAccountsCall: mockListAccountsCall{
				calls:  1,
				filter: model.ListAccountFilter{IsClosed: lo.ToPtr(false)},
				accountList: []model.Account{
					{ID: "1", WalletID: 1, Address: "addr1"},
					{ID: "2", WalletID: 2, Address: "addr2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			dbPort := portsmocks.NewMockDatabasePort(t)

			if tt.mockListAccountsCall.calls > 0 {
				param := tt.mockListAccountsCall
				dbPort.On("ListAccounts", ctx, param.filter).Return(param.accountList, param.expectError).Times(param.calls)
			}

			txPort := portsmocks.NewMockDatabaseTransactionPort(t)

			transactionPort := New(ctx, &Options{
				DatabasePort:    dbPort,
				TransactionPort: dbPort,
				TxPort:          txPort,
				Interval:        1 * time.Minute,
			})

			require.Equal(t, tt.accountList, transactionPort.accountList)
			dbPort.AssertExpectations(t)
			txPort.AssertExpectations(t)
		})
	}
}
func TestTransaction_UpdateAccounts(t *testing.T) {
	t.Parallel()

	type mockListAccountsCall struct {
		calls       int
		filter      model.ListAccountFilter
		accountList []model.Account
		expectError error
	}

	tests := []struct {
		name        string
		accountList map[model.Address]*model.Account
		expectError error

		mockListAccountsCall
	}{
		{
			name: "successful update",
			accountList: map[model.Address]*model.Account{
				"addr1": {ID: "1", WalletID: 1, Address: "addr1"},
				"addr2": {ID: "2", WalletID: 2, Address: "addr2"},
			},
			mockListAccountsCall: mockListAccountsCall{
				calls:  1,
				filter: model.ListAccountFilter{IsClosed: lo.ToPtr(false)},
				accountList: []model.Account{
					{ID: "1", WalletID: 1, Address: "addr1"},
					{ID: "2", WalletID: 2, Address: "addr2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// DBPort mock
			dbPort := portsmocks.NewMockDatabasePort(t)
			if tt.mockListAccountsCall.calls > 0 {
				param := tt.mockListAccountsCall
				dbPort.On("ListAccounts", ctx, param.filter).Return(param.accountList, param.expectError).Times(param.calls)
			}

			transaction := &Transaction{
				dbPort:      dbPort,
				txPort:      portsmocks.NewMockDatabaseTransactionPort(t),
				accountList: make(map[model.Address]*model.Account),
				interval:    1 * time.Minute,
			}

			err := transaction.updateAccounts(ctx)

			require.Equal(t, tt.expectError, err)
			require.Equal(t, tt.accountList, transaction.accountList)

			dbPort.AssertExpectations(t)
		})
	}
}

func TestTransaction_Handle(t *testing.T) {
	t.Parallel()

	var testTransactionMsg = []byte(`{
			"AccountAddr": "13x3kbyGcH8Tt242ZfI21IdodH1f2c8vxeRBK0HFR8s=",
			"LT": 32109106000003,
			"PrevTxHash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			"PrevTxLT": 0,
			"Now": 1741335686,
			"OutMsgCount": 0,
			"OrigStatus": "NON_EXIST",
			"EndStatus": "ACTIVE",
			"IO": {
				"In": {
					"MsgType": "INTERNAL",
					"Msg": {
						"IHRDisabled": true,
						"Bounce": true,
						"Bounced": false,
						"SrcAddr": "EQDNoXNKSXRvOzh3lpUcaeiQY1dxzG6wE6uYn_Cwoh80iIMp",
						"DstAddr": "EQDXfHeRvIZwfxO3bjZl8jbUh2h0fV_Zzy_F5EErQcVHyz3R",
						"Amount": "50000000",
						"ExtraCurrencies": {},
						"IHRFee": "0",
						"FwdFee": "782940",
						"CreatedLT": 32109106000002,
						"CreatedAt": 1741335686,
						"StateInit": {
							"Depth": null,
							"TickTock": null,
							"Code": "te6cckEBAwEAwAABFP8A9KQT9LzyyAsBAfjTMzHQgvCaXeitoBM4xPJe/cDZ81EhGXCHIC1HQo65uifAJNAgxHDIygfL/8nQAfpAMCHHBZWBEzfy8N4BxwCRMOCLcxMDAwMDg4j4J28QgghMS0ChcFMAdIAYyMsFywJQBs8WUAP6AhTLassfIc8WEssfIc8WIc8WIc8WAgBkIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WIc8WAc8WyYAQ+wAzGwRb",
							"Data": "te6cckEBAQEACgAAEAAAAAEAAAAA2uaIhw==",
							"Lib": {}
						},
						"Body": "te6cckEBAQEAAgAAAEysuc0="
					}
				},
				"Out": null
			},
			"TotalFees": {
				"Coins": "532800",
				"ExtraCurrencies": {}
			},
			"StateUpdate": {
				"OldHash": "kK7Illr6uxbrw8ubQI665xthjXh4i8gNCYQ1k8rJjaQ=",
				"NewHash": "aDMXTdFEO6KP8u9lcqHiHvTrFDhA3C6NFhG7m7+2hEw="
			},
			"Description": {
				"CreditFirst": false,
				"StoragePhase": {
					"StorageFeesCollected": "0",
					"StorageFeesDue": null,
					"StatusChange": {
						"Type": "UNCHANGED"
					}
				},
				"CreditPhase": {
					"DueFeesCollected": null,
					"Credit": {
						"Coins": "50000000",
						"ExtraCurrencies": {}
					}
				},
				"ComputePhase": {
					"Phase": {
						"Success": true,
						"MsgStateUsed": false,
						"AccountActivated": false,
						"GasFees": "532800",
						"Details": {
							"GasUsed": 1332,
							"GasLimit": 125000,
							"GasCredit": null,
							"Mode": 0,
							"ExitCode": 0,
							"ExitArg": null,
							"VMSteps": 26,
							"VMInitStateHash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
							"VMFinalStateHash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
						}
					}
				},
				"ActionPhase": {
					"Success": true,
					"Valid": true,
					"NoFunds": false,
					"StatusChange": {
						"Type": "UNCHANGED"
					},
					"TotalFwdFees": null,
					"TotalActionFees": null,
					"ResultCode": 0,
					"ResultArg": null,
					"TotalActions": 0,
					"SpecActions": 0,
					"SkippedActions": 0,
					"MessagesCreated": 0,
					"ActionListHash": "lqKW0iTyhcZ77pPDD4owkVfw2qNdxbh+QQt4YwoJz8c=",
					"TotalMsgSize": {
						"Cells": 0,
						"Bits": 0
					}
				},
				"Aborted": false,
				"BouncePhase": null,
				"Destroyed": false
			},
			"Hash": null
		}`)

	type mockInsertTransactionCall struct {
		calls         int
		tx            *model.Transaction
		responseTx    *model.Transaction
		responseError error
	}

	tests := []struct {
		name        string
		message     []byte
		accountList map[model.Address]*model.Account
		expectError error

		mockInsertTransactionCall
	}{
		{
			name:    "successful",
			message: testTransactionMsg,
			accountList: map[model.Address]*model.Account{
				"EQDNoXNKSXRvOzh3lpUcaeiQY1dxzG6wE6uYn_Cwoh80iIMp": {
					ID:       "1",
					WalletID: 1,
					Address:  "EQDNoXNKSXRvOzh3lpUcaeiQY1dxzG6wE6uYn_Cwoh80iIMp",
				},
			},
			mockInsertTransactionCall: mockInsertTransactionCall{
				calls: 1,
				tx: &model.Transaction{
					AccountAddr: "tx1",
					Sender:      "EQDNoXNKSXRvOzh3lpUcaeiQY1dxzG6wE6uYn_Cwoh80iIMp",
					Receiver:    "EQDXfHeRvIZwfxO3bjZl8jbUh2h0fV_Zzy_F5EErQcVHyz3R",
					Amount:      100,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// TxPort mock
			txPort := portsmocks.NewMockDatabaseTransactionPort(t)
			txPort.On("WithInTransaction", ctx, mock.Anything).
				Return(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

			// TransactionPort mock
			transactionPort := portsmocks.NewMockTransactionalDatabasePort(t)
			if tt.mockInsertTransactionCall.calls > 0 {
				param := tt.mockInsertTransactionCall

				// TODO: fix mock.Anything
				transactionPort.On("InsertTransaction", ctx, mock.Anything).
					Return(param.responseTx, param.responseError).Times(param.calls)
			}

			transaction := &Transaction{
				dbPort:      portsmocks.NewMockDatabasePort(t),
				txPort:      txPort,
				transaction: transactionPort,
				accountList: tt.accountList,
				interval:    1 * time.Minute,
			}

			err := transaction.Handle(ctx, tt.message)
			require.ErrorIs(t, tt.expectError, err)

			transactionPort.AssertExpectations(t)
			txPort.AssertExpectations(t)
		})
	}
}
