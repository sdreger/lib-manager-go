// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

//go:build !build

package book

import (
	"context"

	"github.com/sdreger/lib-manager-go/internal/paging"
	mock "github.com/stretchr/testify/mock"
)

// NewMockStore creates a new instance of MockStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStore {
	mock := &MockStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockStore is an autogenerated mock type for the Store type
type MockStore struct {
	mock.Mock
}

type MockStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStore) EXPECT() *MockStore_Expecter {
	return &MockStore_Expecter{mock: &_m.Mock}
}

// GetByID provides a mock function for the type MockStore
func (_mock *MockStore) GetByID(ctx context.Context, bookID int64) (Book, error) {
	ret := _mock.Called(ctx, bookID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 Book
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int64) (Book, error)); ok {
		return returnFunc(ctx, bookID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, int64) Book); ok {
		r0 = returnFunc(ctx, bookID)
	} else {
		r0 = ret.Get(0).(Book)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = returnFunc(ctx, bookID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockStore_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockStore_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx
//   - bookID
func (_e *MockStore_Expecter) GetByID(ctx interface{}, bookID interface{}) *MockStore_GetByID_Call {
	return &MockStore_GetByID_Call{Call: _e.mock.On("GetByID", ctx, bookID)}
}

func (_c *MockStore_GetByID_Call) Run(run func(ctx context.Context, bookID int64)) *MockStore_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockStore_GetByID_Call) Return(book Book, err error) *MockStore_GetByID_Call {
	_c.Call.Return(book, err)
	return _c
}

func (_c *MockStore_GetByID_Call) RunAndReturn(run func(ctx context.Context, bookID int64) (Book, error)) *MockStore_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Lookup provides a mock function for the type MockStore
func (_mock *MockStore) Lookup(ctx context.Context, page paging.PageRequest, sort paging.Sort, filter Filter) ([]LookupItem, int64, error) {
	ret := _mock.Called(ctx, page, sort, filter)

	if len(ret) == 0 {
		panic("no return value specified for Lookup")
	}

	var r0 []LookupItem
	var r1 int64
	var r2 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, paging.PageRequest, paging.Sort, Filter) ([]LookupItem, int64, error)); ok {
		return returnFunc(ctx, page, sort, filter)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, paging.PageRequest, paging.Sort, Filter) []LookupItem); ok {
		r0 = returnFunc(ctx, page, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]LookupItem)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, paging.PageRequest, paging.Sort, Filter) int64); ok {
		r1 = returnFunc(ctx, page, sort, filter)
	} else {
		r1 = ret.Get(1).(int64)
	}
	if returnFunc, ok := ret.Get(2).(func(context.Context, paging.PageRequest, paging.Sort, Filter) error); ok {
		r2 = returnFunc(ctx, page, sort, filter)
	} else {
		r2 = ret.Error(2)
	}
	return r0, r1, r2
}

// MockStore_Lookup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Lookup'
type MockStore_Lookup_Call struct {
	*mock.Call
}

// Lookup is a helper method to define mock.On call
//   - ctx
//   - page
//   - sort
//   - filter
func (_e *MockStore_Expecter) Lookup(ctx interface{}, page interface{}, sort interface{}, filter interface{}) *MockStore_Lookup_Call {
	return &MockStore_Lookup_Call{Call: _e.mock.On("Lookup", ctx, page, sort, filter)}
}

func (_c *MockStore_Lookup_Call) Run(run func(ctx context.Context, page paging.PageRequest, sort paging.Sort, filter Filter)) *MockStore_Lookup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(paging.PageRequest), args[2].(paging.Sort), args[3].(Filter))
	})
	return _c
}

func (_c *MockStore_Lookup_Call) Return(lookupItems []LookupItem, n int64, err error) *MockStore_Lookup_Call {
	_c.Call.Return(lookupItems, n, err)
	return _c
}

func (_c *MockStore_Lookup_Call) RunAndReturn(run func(ctx context.Context, page paging.PageRequest, sort paging.Sort, filter Filter) ([]LookupItem, int64, error)) *MockStore_Lookup_Call {
	_c.Call.Return(run)
	return _c
}
