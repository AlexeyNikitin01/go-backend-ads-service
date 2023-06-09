// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	user "ads/internal/user"

	mock "github.com/stretchr/testify/mock"
)

// RepositoryUser is an autogenerated mock type for the RepositoryUser type
type RepositoryUser struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: ctx, _a1
func (_m *RepositoryUser) AddUser(ctx context.Context, _a1 *user.User) (int64, error) {
	ret := _m.Called(ctx, _a1)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) (int64, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) int64); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.User) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckUser provides a mock function with given fields: ctx, user_id
func (_m *RepositoryUser) CheckUser(ctx context.Context, user_id int64) bool {
	ret := _m.Called(ctx, user_id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(ctx, user_id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: ctx, user_id
func (_m *RepositoryUser) DeleteUser(ctx context.Context, user_id int64) error {
	ret := _m.Called(ctx, user_id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, user_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx, user_id
func (_m *RepositoryUser) GetUser(ctx context.Context, user_id int64) (*user.User, error) {
	ret := _m.Called(ctx, user_id)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*user.User, error)); ok {
		return rf(ctx, user_id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *user.User); ok {
		r0 = rf(ctx, user_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, user_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, nickname, email, userID, activate
func (_m *RepositoryUser) UpdateUser(ctx context.Context, nickname string, email string, userID int64, activate bool) (*user.User, error) {
	ret := _m.Called(ctx, nickname, email, userID, activate)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64, bool) (*user.User, error)); ok {
		return rf(ctx, nickname, email, userID, activate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64, bool) *user.User); ok {
		r0 = rf(ctx, nickname, email, userID, activate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int64, bool) error); ok {
		r1 = rf(ctx, nickname, email, userID, activate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepositoryUser interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepositoryUser creates a new instance of RepositoryUser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositoryUser(t mockConstructorTestingTNewRepositoryUser) *RepositoryUser {
	mock := &RepositoryUser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
