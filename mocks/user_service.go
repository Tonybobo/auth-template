package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tonybobo/auth-template/models"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindUserById(id string) (*models.DBResponse, error) {
	ret := m.Called(id)
	var r0 *models.DBResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.DBResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
	ret := m.Called(field, value)
	var r0 *models.DBResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.DBResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) ForgetPassword(email string) *models.AuthServiceResponse {
	ret := m.Called(email)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}

func (m *MockUserService) RefreshAccessToken(cookie string) *models.AuthServiceResponse {
	ret := m.Called(cookie)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}

func (m *MockUserService) ResetPassword(user *models.ResetPasswordInput, resetToken string) *models.AuthServiceResponse {
	ret := m.Called(user, resetToken)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}

func (m *MockUserService) VerifyEmail(verificationCode string) *models.AuthServiceResponse {
	ret := m.Called(verificationCode)
	var r0 *models.AuthServiceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.AuthServiceResponse)
	}

	return r0
}
