package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tonybobo/auth-template/models"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Test() *models.AuthServiceResponse {
	ret := m.Called()
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}

func (m *MockAuthService) SignUpUser(user *models.SignUpInput) *models.AuthServiceResponse {
	ret := m.Called(user)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}

func (m *MockAuthService) SignInUser(user *models.SignInInput) *models.AuthServiceResponse {
	ret := m.Called(user)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}
