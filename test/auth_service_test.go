package test

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
	"github.com/tonybobo/auth-template/services"
)

func TestSignUp(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := services.NewAuthService(mockAuthRepository, ctx, temp)
	t.Run("Success", func(t *testing.T) {
		mockUser := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuangjie@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}

		mockUserResp := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bochuang@gmail.com",
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
			Role:     "user",
			Verified: false,
		}

		mockResponse := &models.AuthServiceResponse{
			User:       mockUserResp,
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "An email with the verification code has been sent to " + mockUserResp.Email,
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockAuthRepository.On("SignUpUser", mockArgs1...).Return(mockUserResp, "12345678", nil)

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

		mockResponse := &models.AuthServiceResponse{
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

		mockResponse := &models.AuthServiceResponse{
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

		mockAuthRepository.On("SignUpUser", mockArgs1...).Return(&models.DBResponse{}, "", errors.New("user with that email already exist"))

		response := us.SignUpUser(mockUser)
		assert.Error(t, mockResponse.Err, response.Err)
		assert.Equal(t, mockResponse, response)
		mockAuthRepository.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := services.NewAuthService(mockAuthRepository, ctx, temp)

	t.Run("Success", func(t *testing.T) {
		mockUser := &models.SignInInput{
			Email:    "bochuangjie@gmail.com",
			Password: "12345678",
		}

		mockUserResp := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bochuangjie@gmail.com",
			Role:     "user",
			Verified: true,
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
		}

		mockResponse := &models.AuthServiceResponse{
			User:       mockUserResp,
			Status:     "success",
			StatusCode: http.StatusOK,
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.Email,
		}

		mockAuthRepository.On("FindUserByEmail", mockArgs1...).Return(mockUserResp, nil)

		response := us.SignInUser(mockUser)
		assert.NoError(t, response.Err)
		assert.Equal(t, response.Status, mockResponse.Status)
		assert.Equal(t, response.StatusCode, mockResponse.StatusCode)
		mockAuthRepository.AssertExpectations(t)
	})

	t.Run("User not verified", func(t *testing.T) {
		mockUser := &models.SignInInput{
			Email:    "bobo@gmail.com",
			Password: "12345678",
		}

		mockAuthResponse := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bobo@gmail.com",
			Role:     "user",
			Verified: false,
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
		}

		mockResponse := &models.AuthServiceResponse{
			User:       nil,
			Status:     "fail",
			StatusCode: http.StatusUnauthorized,
			Message:    "You have not verify the account , Please verify your email to login",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.Email,
		}

		mockAuthRepository.On("FindUserByEmail", mockArgs...).Return(mockAuthResponse, nil)

		response := us.SignInUser(mockUser)
		assert.Error(t, response.Err)
		assert.Equal(t, response.Status, mockResponse.Status)
		assert.Equal(t, response.StatusCode, mockResponse.StatusCode)
		mockAuthRepository.AssertExpectations(t)
	})
	t.Run("Incorrect Password", func(t *testing.T) {
		mockUser := &models.SignInInput{
			Email:    "bobo@gmail.com",
			Password: "1234567",
		}

		mockAuthResponse := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bobo@gmail.com",
			Role:     "user",
			Verified: true,
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
		}

		mockResponse := &models.AuthServiceResponse{
			User:       nil,
			Status:     "fail",
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid Email or Password",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.Email,
		}

		mockAuthRepository.On("FindUserByEmail", mockArgs...).Return(mockAuthResponse, nil)

		response := us.SignInUser(mockUser)
		assert.Error(t, response.Err)
		assert.Equal(t, response.Status, mockResponse.Status)
		assert.Equal(t, response.StatusCode, mockResponse.StatusCode)
		mockAuthRepository.AssertExpectations(t)

	})
}
