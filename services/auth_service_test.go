package services

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tonybobo/auth-template/mocks"
	"github.com/tonybobo/auth-template/models"
)

func TestSignUp(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := NewAuthService(mockAuthRepository, ctx, temp)
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuangjie@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}

		mockUserResp := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bochuangjie@gmail.com",
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
			Role:     "user",
			Verified: false,
		}

		mockResponse := &AuthServiceResponse{
			User:       mockUserResp,
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "An email with the verification code has been sent to " + mockUserResp.Email,
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockAuthRepository.On("SignUpUser", mockArgs1...).Return(mockUserResp, nil, "12345678")

		response := us.SignUpUser(mockUser)
		assert.NoError(t, response.Err)
		assert.Equal(t, response.Status, mockResponse.Status)
		assert.Equal(t, response.StatusCode, mockResponse.StatusCode)
		mockAuthRepository.AssertExpectations(t)
	})

	t.Run("Password not match", func(t *testing.T) {
		mockUser := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuangjie@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "123456",
		}

		mockResponse := &AuthServiceResponse{
			User:       nil,
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
			Message:    "password not match",
			Err:        errors.New("password not match"),
		}

		response := us.SignUpUser(mockUser)
		assert.Error(t, mockResponse.Err, response.Err)
		assert.Equal(t, mockResponse, response)
		mockAuthRepository.AssertNotCalled(t, "SignUpUser")
	})

	t.Run("User already exist", func(t *testing.T) {
		mockUser := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuangjie@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}

		mockResponse := &AuthServiceResponse{
			User:       nil,
			Status:     "fail",
			StatusCode: http.StatusBadGateway,
			Message:    "user with that email already exist",
			Err:        errors.New("user with that email already exist"),
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockAuthRepository.On("SignUpUser", mockArgs1...).Return(&models.DBResponse{}, errors.New("user with that email already exist"), "")

		response := us.SignUpUser(mockUser)
		assert.Error(t, mockResponse.Err, response.Err)
		assert.Equal(t, mockResponse, response)
		mockAuthRepository.AssertExpectations(t)
	})
}
