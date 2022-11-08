package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tonybobo/auth-template/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) SignUpUser(ctx context.Context, user *models.SignUpInput) (*models.DBResponse, string, error) {
	ret := m.Called(ctx, user)

	var r0 *models.DBResponse
	var r1 string
	var r2 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.DBResponse)
	}

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(string)
	}
	if ret.Get(2) != nil {
		r2 = ret.Get(2).(error)
	}

	return r0, r1, r2
}

func (m *MockAuthRepository) ClearResetPasswordToken(ctx context.Context, token, password string) error {
	ret := m.Called(ctx, token, password)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockAuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.DBResponse, error) {
	ret := m.Called(ctx, email)

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

func (m *MockAuthRepository) FindUserById(ctx context.Context, id primitive.ObjectID) (*models.DBResponse, error) {
	ret := m.Called(ctx, id)

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

func (m *MockAuthRepository) ResetPasswordToken(ctx context.Context, email, passwordResetToken string) (*mongo.UpdateResult, error) {
	ret := m.Called(ctx, email, passwordResetToken)

	var r0 *mongo.UpdateResult

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*mongo.UpdateResult)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockAuthRepository) UpdateOne(ctx context.Context, field string, value interface{}) (*mongo.UpdateResult, error) {
	ret := m.Called(ctx, field, value)

	var r0 *mongo.UpdateResult

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*mongo.UpdateResult)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockAuthRepository) VerifyEmail(ctx context.Context, verificationCode string) error {
	ret := m.Called(ctx, verificationCode)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockAuthRepository) ForgetPassword(ctx context.Context, email string) (*models.DBResponse, string, error) {
	ret := m.Called(ctx, email)

	var r0 *models.DBResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.DBResponse)
	}

	var r1 string

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(string)
	}

	var r2 error

	if ret.Get(2) != nil {
		r2 = ret.Get(2).(error)
	}

	return r0, r1, r2
}
