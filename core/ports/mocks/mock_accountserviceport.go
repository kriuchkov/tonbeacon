// Code generated by mockery v2.52.3. DO NOT EDIT.

package portsmocks

import (
	context "context"

	model "github.com/kriuchkov/tonbeacon/core/model"
	mock "github.com/stretchr/testify/mock"
)

// MockAccountServicePort is an autogenerated mock type for the AccountServicePort type
type MockAccountServicePort struct {
	mock.Mock
}

type MockAccountServicePort_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccountServicePort) EXPECT() *MockAccountServicePort_Expecter {
	return &MockAccountServicePort_Expecter{mock: &_m.Mock}
}

// CloseAccount provides a mock function with given fields: ctx, accountID
func (_m *MockAccountServicePort) CloseAccount(ctx context.Context, accountID string) error {
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

// MockAccountServicePort_CloseAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CloseAccount'
type MockAccountServicePort_CloseAccount_Call struct {
	*mock.Call
}

// CloseAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockAccountServicePort_Expecter) CloseAccount(ctx interface{}, accountID interface{}) *MockAccountServicePort_CloseAccount_Call {
	return &MockAccountServicePort_CloseAccount_Call{Call: _e.mock.On("CloseAccount", ctx, accountID)}
}

func (_c *MockAccountServicePort_CloseAccount_Call) Run(run func(ctx context.Context, accountID string)) *MockAccountServicePort_CloseAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAccountServicePort_CloseAccount_Call) Return(_a0 error) *MockAccountServicePort_CloseAccount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccountServicePort_CloseAccount_Call) RunAndReturn(run func(context.Context, string) error) *MockAccountServicePort_CloseAccount_Call {
	_c.Call.Return(run)
	return _c
}

// CreateAccount provides a mock function with given fields: ctx, accountID
func (_m *MockAccountServicePort) CreateAccount(ctx context.Context, accountID string) (*model.Account, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for CreateAccount")
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

// MockAccountServicePort_CreateAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAccount'
type MockAccountServicePort_CreateAccount_Call struct {
	*mock.Call
}

// CreateAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockAccountServicePort_Expecter) CreateAccount(ctx interface{}, accountID interface{}) *MockAccountServicePort_CreateAccount_Call {
	return &MockAccountServicePort_CreateAccount_Call{Call: _e.mock.On("CreateAccount", ctx, accountID)}
}

func (_c *MockAccountServicePort_CreateAccount_Call) Run(run func(ctx context.Context, accountID string)) *MockAccountServicePort_CreateAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAccountServicePort_CreateAccount_Call) Return(_a0 *model.Account, _a1 error) *MockAccountServicePort_CreateAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountServicePort_CreateAccount_Call) RunAndReturn(run func(context.Context, string) (*model.Account, error)) *MockAccountServicePort_CreateAccount_Call {
	_c.Call.Return(run)
	return _c
}

// GetBalance provides a mock function with given fields: ctx, accountID
func (_m *MockAccountServicePort) GetBalance(ctx context.Context, accountID string) ([]model.Balance, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 []model.Balance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]model.Balance, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Balance); ok {
		r0 = rf(ctx, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Balance)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountServicePort_GetBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBalance'
type MockAccountServicePort_GetBalance_Call struct {
	*mock.Call
}

// GetBalance is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID string
func (_e *MockAccountServicePort_Expecter) GetBalance(ctx interface{}, accountID interface{}) *MockAccountServicePort_GetBalance_Call {
	return &MockAccountServicePort_GetBalance_Call{Call: _e.mock.On("GetBalance", ctx, accountID)}
}

func (_c *MockAccountServicePort_GetBalance_Call) Run(run func(ctx context.Context, accountID string)) *MockAccountServicePort_GetBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAccountServicePort_GetBalance_Call) Return(_a0 []model.Balance, _a1 error) *MockAccountServicePort_GetBalance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountServicePort_GetBalance_Call) RunAndReturn(run func(context.Context, string) ([]model.Balance, error)) *MockAccountServicePort_GetBalance_Call {
	_c.Call.Return(run)
	return _c
}

// ListAccounts provides a mock function with given fields: ctx, req
func (_m *MockAccountServicePort) ListAccounts(ctx context.Context, req model.ListAccountFilter) ([]model.Account, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ListAccounts")
	}

	var r0 []model.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ListAccountFilter) ([]model.Account, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ListAccountFilter) []model.Account); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ListAccountFilter) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountServicePort_ListAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAccounts'
type MockAccountServicePort_ListAccounts_Call struct {
	*mock.Call
}

// ListAccounts is a helper method to define mock.On call
//   - ctx context.Context
//   - req model.ListAccountFilter
func (_e *MockAccountServicePort_Expecter) ListAccounts(ctx interface{}, req interface{}) *MockAccountServicePort_ListAccounts_Call {
	return &MockAccountServicePort_ListAccounts_Call{Call: _e.mock.On("ListAccounts", ctx, req)}
}

func (_c *MockAccountServicePort_ListAccounts_Call) Run(run func(ctx context.Context, req model.ListAccountFilter)) *MockAccountServicePort_ListAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.ListAccountFilter))
	})
	return _c
}

func (_c *MockAccountServicePort_ListAccounts_Call) Return(_a0 []model.Account, _a1 error) *MockAccountServicePort_ListAccounts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountServicePort_ListAccounts_Call) RunAndReturn(run func(context.Context, model.ListAccountFilter) ([]model.Account, error)) *MockAccountServicePort_ListAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// MasterAccount provides a mock function with given fields: ctx
func (_m *MockAccountServicePort) MasterAccount(ctx context.Context) (*model.Account, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for MasterAccount")
	}

	var r0 *model.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*model.Account, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *model.Account); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountServicePort_MasterAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MasterAccount'
type MockAccountServicePort_MasterAccount_Call struct {
	*mock.Call
}

// MasterAccount is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAccountServicePort_Expecter) MasterAccount(ctx interface{}) *MockAccountServicePort_MasterAccount_Call {
	return &MockAccountServicePort_MasterAccount_Call{Call: _e.mock.On("MasterAccount", ctx)}
}

func (_c *MockAccountServicePort_MasterAccount_Call) Run(run func(ctx context.Context)) *MockAccountServicePort_MasterAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAccountServicePort_MasterAccount_Call) Return(_a0 *model.Account, _a1 error) *MockAccountServicePort_MasterAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountServicePort_MasterAccount_Call) RunAndReturn(run func(context.Context) (*model.Account, error)) *MockAccountServicePort_MasterAccount_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccountServicePort creates a new instance of MockAccountServicePort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountServicePort(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountServicePort {
	mock := &MockAccountServicePort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
