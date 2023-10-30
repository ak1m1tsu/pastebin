// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	entity "github.com/romankravchuk/pastebin/internal/entity"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// AuthWebAPI is an autogenerated mock type for the AuthWebAPI type
type AuthWebAPI struct {
	mock.Mock
}

// GetToken provides a mock function with given fields: code
func (_m *AuthWebAPI) GetToken(code string) (*oauth2.Token, error) {
	ret := _m.Called(code)

	var r0 *oauth2.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*oauth2.Token, error)); ok {
		return rf(code)
	}
	if rf, ok := ret.Get(0).(func(string) *oauth2.Token); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInfo provides a mock function with given fields: token
func (_m *AuthWebAPI) GetUserInfo(token *oauth2.Token) (*entity.APIUser, error) {
	ret := _m.Called(token)

	var r0 *entity.APIUser
	var r1 error
	if rf, ok := ret.Get(0).(func(*oauth2.Token) (*entity.APIUser, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(*oauth2.Token) *entity.APIUser); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.APIUser)
		}
	}

	if rf, ok := ret.Get(1).(func(*oauth2.Token) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthWebAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthWebAPI creates a new instance of AuthWebAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthWebAPI(t mockConstructorTestingTNewAuthWebAPI) *AuthWebAPI {
	mock := &AuthWebAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}