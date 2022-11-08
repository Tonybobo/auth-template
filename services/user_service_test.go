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

func TestRefreshAccessToken(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := NewUserServiceImpl(mockAuthRepository, ctx, temp)

	t.Run("expired token", func(t *testing.T) {

		cookie := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njc4MTM1ODIsImlhdCI6MTY2NzgwOTk4MiwibmJmIjoxNjY3ODA5OTgyLCJzdWIiOiI2MzY3NjkzZGRlMjkwOTFhYWFhNWRhZGQifQ.Xt3UtH6ll9q1y0vVelCU0cgPt_gKTKHLigIKHggOW1IbEFMvhtsSHI_Sh7dJlq73F1Sgs8KTleX19KoWOQb3bw"
		mockResponse := &AuthServiceResponse{
			User:       nil,
			Status:     "fail",
			StatusCode: http.StatusForbidden,
		}

		response := us.RefreshAccessToken(cookie)
		assert.Error(t, response.Err)
		assert.Equal(t, mockResponse.Status, response.Status)
		assert.Equal(t, mockResponse.StatusCode, response.StatusCode)

	})

}

func TestVerifyEmail(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := NewUserServiceImpl(mockAuthRepository, ctx, temp)

	t.Run("Success", func(t *testing.T) {

		mockResponse := &AuthServiceResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "Successfully Verified",
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"asdasdasddaa",
		}
		mockAuthRepository.On("VerifyEmail", mockArg...).Return(nil)

		response := us.VerifyEmail("asdasdasddaa")
		assert.NoError(t, response.Err)
		assert.ObjectsAreEqualValues(mockResponse, response)
	})

	t.Run("Invalid Email", func(t *testing.T) {

		mockResponse := &AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusForbidden,
			Message:    "invalid email",
			Err:        errors.New("invalid email"),
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"12345678",
		}
		mockAuthRepository.On("VerifyEmail", mockArg...).Return(errors.New("invalid email"))

		response := us.VerifyEmail("12345678")
		assert.Error(t, response.Err)
		assert.ObjectsAreEqualValues(mockResponse, response)
	})

	t.Run("DB error", func(t *testing.T) {

		mockResponse := &AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusBadGateway,
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"22345678",
		}
		mockAuthRepository.On("VerifyEmail", mockArg...).Return(errors.New("db error"))

		response := us.VerifyEmail("22345678")
		assert.Error(t, response.Err)
		assert.Equal(t, mockResponse.Status, response.Status)
		assert.Equal(t, mockResponse.StatusCode, response.StatusCode)
	})

}

func TestResetPassword(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := NewUserServiceImpl(mockAuthRepository, ctx, temp)

	t.Run("Success", func(t *testing.T) {
		mockUserInput := &models.ResetPasswordInput{
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}

		mockResponse := &AuthServiceResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "Password updated successfully. Please Login with new password",
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"YXNkYXNzZGFzZGFkYXNkYXM=",
			mockUserInput.Password,
		}
		mockAuthRepository.On("ClearResetPasswordToken", mockArg...).Return(nil)

		response := us.ResetPassword(mockUserInput, "YXNkYXNzZGFzZGFkYXNkYXM=")
		assert.NoError(t, response.Err)
		assert.ObjectsAreEqualValues(mockResponse, response)
	})

	t.Run("invalid or expired token", func(t *testing.T) {
		mockUserInput := &models.ResetPasswordInput{
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}

		mockResponse := &AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusForbidden,
			Message:    "invalid or expired token",
			Err:        errors.New("invalid or expired token"),
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"YXNkYXNzZGFzZGFkYXNkYXM=",
			mockUserInput.Password,
		}
		mockAuthRepository.On("ClearResetPasswordToken", mockArg...).Return(errors.New("invalid or expired token"))

		response := us.ResetPassword(mockUserInput, "YXNkYXNzZGFzZGFkYXNkYXM=")

		assert.ObjectsAreEqualValues(mockResponse, response)
	})

	t.Run("Password not match", func(t *testing.T) {
		mockUserInput := &models.ResetPasswordInput{
			Password:        "12345678",
			PasswordConfirm: "12345679",
		}

		mockResponse := &AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
			Message:    "Password does not match",
		}
		mockArg := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			"YXNkYXNzZGFzZGFkYXNkYXM",
			mockUserInput.Password,
		}

		response := us.ResetPassword(mockUserInput, "YXNkYXNzZGFzZGFkYXNkYXM")

		assert.ObjectsAreEqualValues(mockResponse, response)
		mockAuthRepository.AssertNotCalled(t, "ClearResetPasswordToken", mockArg...)
	})

}

func TestForgetPassword(t *testing.T) {
	mockAuthRepository := new(mocks.MockAuthRepository)
	ctx := context.TODO()
	temp := template.Must(template.ParseGlob("../templates/*.html"))
	us := NewUserServiceImpl(mockAuthRepository, ctx, temp)

	t.Run("Success", func(t *testing.T) {
		email := "bochuang@gmail.com"
		resetToken := "123456"
		mockResponse := &AuthServiceResponse{
			Message:    "You will receive a reset email if user with that email exist",
			Status:     "success",
			StatusCode: http.StatusOK,
		}

		mockUserResp := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bochuang@gmail.com",
			Role:     "user",
			Verified: true,
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockAuthRepository.On("ForgetPassword", mockArgs1...).Return(mockUserResp, resetToken, nil)

		response := us.ForgetPassword(email)
		assert.NoError(t, response.Err)
		assert.ObjectsAreEqual(mockResponse, response)
	})

	t.Run("On User Not Verified", func(t *testing.T) {
		email := "bobo@gmail.com"
		resetToken := "12345"
		mockResponse := &AuthServiceResponse{
			Message:    "account has not been verified. please verify your account with the email sent",
			Status:     "fail",
			StatusCode: http.StatusUnauthorized,
		}

		mockUserResp := &models.DBResponse{
			Name:     "Bo Chuang Jie",
			Email:    "bobo@gmail.com",
			Role:     "user",
			Verified: false,
			Password: "$2a$10$AxVZoSa4XI1XdTbvElmm2eNHsBm7KST02qTmGboWYOleB4NOv11PK",
		}

		mockArgs1 := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockAuthRepository.On("ForgetPassword", mockArgs1...).Return(mockUserResp, resetToken, nil)

		response := us.ForgetPassword(email)
		assert.Error(t, response.Err)
		assert.ObjectsAreEqual(mockResponse, response)
	})
}
