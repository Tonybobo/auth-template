package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tonybobo/auth-template/models"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) FindUserById(id string) (*models.DBResponse, error) {
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
