// Code generated by mockery v2.52.3. DO NOT EDIT.

package portsmocks

import (
	context "context"

	model "github.com/kriuchkov/tonbeacon/core/model"
	mock "github.com/stretchr/testify/mock"
)

// MockDatabasePort is an autogenerated mock type for the DatabasePort type
type MockDatabasePort struct {
	mock.Mock
}

type MockDatabasePort_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabasePort) EXPECT() *MockDatabasePort_Expecter {
	return &MockDatabasePort_Expecter{mock: &_m.Mock}
}

// CloseAccount provides a mock function with given fields: ctx, accountID
func (_m *MockDatabasePort) CloseAccount(ctx context.Context, accountID string) error {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for CloseAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, accountID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabasePort_CloseAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CloseAccount'
type MockDatabasePort_CloseAccount_Call struct {
	*mock.Call
}

// CloseAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockDatabasePort_Expecter) CloseAccount(ctx interface{}, accountID interface{}) *MockDatabasePort_CloseAccount_Call {
	return &MockDatabasePort_CloseAccount_Call{Call: _e.mock.On("CloseAccount", ctx, accountID)}
}

func (_c *MockDatabasePort_CloseAccount_Call) Run(run func(ctx context.Context, accountID string)) *MockDatabasePort_CloseAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabasePort_CloseAccount_Call) Return(_a0 error) *MockDatabasePort_CloseAccount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabasePort_CloseAccount_Call) RunAndReturn(run func(context.Context, string) error) *MockDatabasePort_CloseAccount_Call {
	_c.Call.Return(run)
	return _c
}

// GetEvents provides a mock function with given fields: ctx, limit
func (_m *MockDatabasePort) GetEvents(ctx context.Context, limit int64) ([]model.OutboxEvent, error) {
	ret := _m.Called(ctx, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetEvents")
	}

	var r0 []model.OutboxEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]model.OutboxEvent, error)); ok {
		return rf(ctx, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []model.OutboxEvent); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.OutboxEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_GetEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEvents'
type MockDatabasePort_GetEvents_Call struct {
	*mock.Call
}

// GetEvents is a helper method to define mock.On call
//   - ctx context.Context
//   - limit int64
func (_e *MockDatabasePort_Expecter) GetEvents(ctx interface{}, limit interface{}) *MockDatabasePort_GetEvents_Call {
	return &MockDatabasePort_GetEvents_Call{Call: _e.mock.On("GetEvents", ctx, limit)}
}

func (_c *MockDatabasePort_GetEvents_Call) Run(run func(ctx context.Context, limit int64)) *MockDatabasePort_GetEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockDatabasePort_GetEvents_Call) Return(_a0 []model.OutboxEvent, _a1 error) *MockDatabasePort_GetEvents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_GetEvents_Call) RunAndReturn(run func(context.Context, int64) ([]model.OutboxEvent, error)) *MockDatabasePort_GetEvents_Call {
	_c.Call.Return(run)
	return _c
}

// GetWalletIDByAccountID provides a mock function with given fields: ctx, accountID
func (_m *MockDatabasePort) GetWalletIDByAccountID(ctx context.Context, accountID string) (uint32, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletIDByAccountID")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (uint32, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) uint32); ok {
		r0 = rf(ctx, accountID)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_GetWalletIDByAccountID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWalletIDByAccountID'
type MockDatabasePort_GetWalletIDByAccountID_Call struct {
	*mock.Call
}

// GetWalletIDByAccountID is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockDatabasePort_Expecter) GetWalletIDByAccountID(ctx interface{}, accountID interface{}) *MockDatabasePort_GetWalletIDByAccountID_Call {
	return &MockDatabasePort_GetWalletIDByAccountID_Call{Call: _e.mock.On("GetWalletIDByAccountID", ctx, accountID)}
}

func (_c *MockDatabasePort_GetWalletIDByAccountID_Call) Run(run func(ctx context.Context, accountID string)) *MockDatabasePort_GetWalletIDByAccountID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabasePort_GetWalletIDByAccountID_Call) Return(_a0 uint32, _a1 error) *MockDatabasePort_GetWalletIDByAccountID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_GetWalletIDByAccountID_Call) RunAndReturn(run func(context.Context, string) (uint32, error)) *MockDatabasePort_GetWalletIDByAccountID_Call {
	_c.Call.Return(run)
	return _c
}

// InsertAccount provides a mock function with given fields: ctx, accountID
func (_m *MockDatabasePort) InsertAccount(ctx context.Context, accountID string) (*model.Account, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for InsertAccount")
	}

	var r0 *model.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Account, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Account); ok {
		r0 = rf(ctx, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_InsertAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertAccount'
type MockDatabasePort_InsertAccount_Call struct {
	*mock.Call
}

// InsertAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockDatabasePort_Expecter) InsertAccount(ctx interface{}, accountID interface{}) *MockDatabasePort_InsertAccount_Call {
	return &MockDatabasePort_InsertAccount_Call{Call: _e.mock.On("InsertAccount", ctx, accountID)}
}

func (_c *MockDatabasePort_InsertAccount_Call) Run(run func(ctx context.Context, accountID string)) *MockDatabasePort_InsertAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabasePort_InsertAccount_Call) Return(_a0 *model.Account, _a1 error) *MockDatabasePort_InsertAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_InsertAccount_Call) RunAndReturn(run func(context.Context, string) (*model.Account, error)) *MockDatabasePort_InsertAccount_Call {
	_c.Call.Return(run)
	return _c
}

// InsertTransaction provides a mock function with given fields: ctx, tx
func (_m *MockDatabasePort) InsertTransaction(ctx context.Context, tx *model.Transaction) (*model.Transaction, error) {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for InsertTransaction")
	}

	var r0 *model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Transaction) (*model.Transaction, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Transaction) *model.Transaction); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Transaction) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_InsertTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertTransaction'
type MockDatabasePort_InsertTransaction_Call struct {
	*mock.Call
}

// InsertTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - tx *model.Transaction
func (_e *MockDatabasePort_Expecter) InsertTransaction(ctx interface{}, tx interface{}) *MockDatabasePort_InsertTransaction_Call {
	return &MockDatabasePort_InsertTransaction_Call{Call: _e.mock.On("InsertTransaction", ctx, tx)}
}

func (_c *MockDatabasePort_InsertTransaction_Call) Run(run func(ctx context.Context, tx *model.Transaction)) *MockDatabasePort_InsertTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Transaction))
	})
	return _c
}

func (_c *MockDatabasePort_InsertTransaction_Call) Return(_a0 *model.Transaction, _a1 error) *MockDatabasePort_InsertTransaction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_InsertTransaction_Call) RunAndReturn(run func(context.Context, *model.Transaction) (*model.Transaction, error)) *MockDatabasePort_InsertTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// IsAccountExists provides a mock function with given fields: ctx, accountID
func (_m *MockDatabasePort) IsAccountExists(ctx context.Context, accountID string) (bool, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for IsAccountExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, accountID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_IsAccountExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAccountExists'
type MockDatabasePort_IsAccountExists_Call struct {
	*mock.Call
}

// IsAccountExists is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockDatabasePort_Expecter) IsAccountExists(ctx interface{}, accountID interface{}) *MockDatabasePort_IsAccountExists_Call {
	return &MockDatabasePort_IsAccountExists_Call{Call: _e.mock.On("IsAccountExists", ctx, accountID)}
}

func (_c *MockDatabasePort_IsAccountExists_Call) Run(run func(ctx context.Context, accountID string)) *MockDatabasePort_IsAccountExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabasePort_IsAccountExists_Call) Return(_a0 bool, _a1 error) *MockDatabasePort_IsAccountExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_IsAccountExists_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *MockDatabasePort_IsAccountExists_Call {
	_c.Call.Return(run)
	return _c
}

// ListAccounts provides a mock function with given fields: ctx, filter
func (_m *MockDatabasePort) ListAccounts(ctx context.Context, filter model.ListAccountFilter) ([]model.Account, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListAccounts")
	}

	var r0 []model.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ListAccountFilter) ([]model.Account, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ListAccountFilter) []model.Account); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ListAccountFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabasePort_ListAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAccounts'
type MockDatabasePort_ListAccounts_Call struct {
	*mock.Call
}

// ListAccounts is a helper method to define mock.On call
//   - ctx context.Context
//   - filter model.ListAccountFilter
func (_e *MockDatabasePort_Expecter) ListAccounts(ctx interface{}, filter interface{}) *MockDatabasePort_ListAccounts_Call {
	return &MockDatabasePort_ListAccounts_Call{Call: _e.mock.On("ListAccounts", ctx, filter)}
}

func (_c *MockDatabasePort_ListAccounts_Call) Run(run func(ctx context.Context, filter model.ListAccountFilter)) *MockDatabasePort_ListAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ListAccountFilter))
	})
	return _c
}

func (_c *MockDatabasePort_ListAccounts_Call) Return(_a0 []model.Account, _a1 error) *MockDatabasePort_ListAccounts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabasePort_ListAccounts_Call) RunAndReturn(run func(context.Context, model.ListAccountFilter) ([]model.Account, error)) *MockDatabasePort_ListAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// MarkEventAsProcessed provides a mock function with given fields: ctx, eventID
func (_m *MockDatabasePort) MarkEventAsProcessed(ctx context.Context, eventID uint64) error {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for MarkEventAsProcessed")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabasePort_MarkEventAsProcessed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarkEventAsProcessed'
type MockDatabasePort_MarkEventAsProcessed_Call struct {
	*mock.Call
}

// MarkEventAsProcessed is a helper method to define mock.On call
//   - ctx context.Context
//   - eventID uint64
func (_e *MockDatabasePort_Expecter) MarkEventAsProcessed(ctx interface{}, eventID interface{}) *MockDatabasePort_MarkEventAsProcessed_Call {
	return &MockDatabasePort_MarkEventAsProcessed_Call{Call: _e.mock.On("MarkEventAsProcessed", ctx, eventID)}
}

func (_c *MockDatabasePort_MarkEventAsProcessed_Call) Run(run func(ctx context.Context, eventID uint64)) *MockDatabasePort_MarkEventAsProcessed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *MockDatabasePort_MarkEventAsProcessed_Call) Return(_a0 error) *MockDatabasePort_MarkEventAsProcessed_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabasePort_MarkEventAsProcessed_Call) RunAndReturn(run func(context.Context, uint64) error) *MockDatabasePort_MarkEventAsProcessed_Call {
	_c.Call.Return(run)
	return _c
}

// SaveEvent provides a mock function with given fields: ctx, event
func (_m *MockDatabasePort) SaveEvent(ctx context.Context, event model.OutboxEvent) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for SaveEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.OutboxEvent) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabasePort_SaveEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveEvent'
type MockDatabasePort_SaveEvent_Call struct {
	*mock.Call
}

// SaveEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - event model.OutboxEvent
func (_e *MockDatabasePort_Expecter) SaveEvent(ctx interface{}, event interface{}) *MockDatabasePort_SaveEvent_Call {
	return &MockDatabasePort_SaveEvent_Call{Call: _e.mock.On("SaveEvent", ctx, event)}
}

func (_c *MockDatabasePort_SaveEvent_Call) Run(run func(ctx context.Context, event model.OutboxEvent)) *MockDatabasePort_SaveEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.OutboxEvent))
	})
	return _c
}

func (_c *MockDatabasePort_SaveEvent_Call) Return(_a0 error) *MockDatabasePort_SaveEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabasePort_SaveEvent_Call) RunAndReturn(run func(context.Context, model.OutboxEvent) error) *MockDatabasePort_SaveEvent_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAccount provides a mock function with given fields: ctx, account
func (_m *MockDatabasePort) UpdateAccount(ctx context.Context, account *model.Account) error {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabasePort_UpdateAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAccount'
type MockDatabasePort_UpdateAccount_Call struct {
	*mock.Call
}

// UpdateAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - account *model.Account
func (_e *MockDatabasePort_Expecter) UpdateAccount(ctx interface{}, account interface{}) *MockDatabasePort_UpdateAccount_Call {
	return &MockDatabasePort_UpdateAccount_Call{Call: _e.mock.On("UpdateAccount", ctx, account)}
}

func (_c *MockDatabasePort_UpdateAccount_Call) Run(run func(ctx context.Context, account *model.Account)) *MockDatabasePort_UpdateAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Account))
	})
	return _c
}

func (_c *MockDatabasePort_UpdateAccount_Call) Return(_a0 error) *MockDatabasePort_UpdateAccount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabasePort_UpdateAccount_Call) RunAndReturn(run func(context.Context, *model.Account) error) *MockDatabasePort_UpdateAccount_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabasePort creates a new instance of MockDatabasePort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabasePort(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabasePort {
	mock := &MockDatabasePort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
