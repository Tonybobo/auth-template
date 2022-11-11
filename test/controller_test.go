package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tonybobo/auth-template/controllers"
	"github.com/tonybobo/auth-template/mocks"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/routes"
	"github.com/tonybobo/auth-template/utils"
)

var (
	mockAuthService     = new(mocks.MockAuthService)
	mockUserService     = new(mocks.MockUserService)
	ctx                 = context.TODO()
	authController      = controllers.NewAuthController(mockAuthService, mockUserService, ctx)
	authRouteController = routes.NewAuthRouteController(authController)
	userController      = controllers.NewUserController(mockUserService)
	userRouteController = routes.NewUserRouteController(userController)
	server              = gin.Default()
	router              = server.Group("/api")
)

func TestAuth(t *testing.T) {
	mockAuthService.On("Test").Return(&models.AuthServiceResponse{
		Message: "test",
	})
	authRouteController.AuthRoute(router, mockUserService)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/auth/test", nil)

	if err != nil {
		t.FailNow()
	}
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":0}", w.Body.String())
}

func TestSignUpUser(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		user := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuang@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "12345678",
		}
		mockResp := &models.AuthServiceResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "An email with the verification code has been sent to " + user.Email,
		}

		reqBody, err := json.Marshal(gin.H{
			"name":            user.Name,
			"email":           user.Email,
			"password":        user.Password,
			"passwordConfirm": user.PasswordConfirm,
		})

		assert.NoError(t, err)

		mockAuthService.On("SignUpUser", user).Return(mockResp)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})

	t.Run("Password Not match", func(t *testing.T) {
		user := &models.SignUpInput{
			Name:            "Bo Chuang Jie",
			Email:           "bochuang@gmail.com",
			Password:        "12345678",
			PasswordConfirm: "1234567",
		}
		mockResp := &models.AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
			Message:    "password not match",
			Err:        errors.New("password not match"),
		}

		reqBody, err := json.Marshal(gin.H{
			"name":            user.Name,
			"email":           user.Email,
			"password":        user.Password,
			"passwordConfirm": user.PasswordConfirm,
		})

		assert.NoError(t, err)

		mockAuthService.On("SignUpUser", user).Return(mockResp)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(reqBody))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})

}

func TestRefreshToken(t *testing.T) {
	t.Run("refresh Token sucessfully", func(t *testing.T) {
		cookie := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjgxNTUwMDAsImlhdCI6MTY2ODE1MTQwMCwibmJmIjoxNjY4MTUxNDAwLCJzdWIiOiI2MzZiN2QwMzFjZjY2ZWFkNmNiZGU5OWEifQ.JXWp7yA5JPRjwXakLkzxSyPXFUujSFYL8_BVZcvGgKyrOYTUEu23DwjWMfyjiEavMk1LMQoJ3GWP6RIgq-07mg"

		mockResp := &models.AuthServiceResponse{
			Status:      "success",
			StatusCode:  http.StatusOK,
			AccessToken: "testing",
		}

		mockUserService.On("RefreshAccessToken", cookie).Return(mockResp)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: cookie})

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":       mockResp.Status,
			"access_token": mockResp.AccessToken,
		})

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)

	})

	t.Run("no cookie in request", func(t *testing.T) {
		mockResp := &models.AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusForbidden,
			Message:    "could not refresh access token",
		}

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

	})

	t.Run("invalid cookie", func(t *testing.T) {
		cookie := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjgxNTUwMDAslhdCI6MTY2ODE1MTQwMCwibmJmIjoxNjY4MTUxNDAwLCJzdWIiOiI2MzZiN2QwMzFjZjY2ZWFkNmNiZGU5OWEifQ.JXWp7yA5JPRjwXakLkzxSyPXFUujSFYL8_BVZcvGgKyrOYTUEu23DwjWMfyjiEavMk1LMQoJ3GWP6RIgq-07mg"

		mockResp := &models.AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusForbidden,
			Message:    "Invalid Token",
			Err:        errors.New("invalid token"),
		}

		mockUserService.On("RefreshAccessToken", cookie).Return(mockResp)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: cookie})

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResp.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)

	})
}

func TestVerifyEmailController(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		code := "1234423"
		verificationCode := utils.Encode(code)
		mockResp := &models.AuthServiceResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "Successfully Verified",
		}

		mockUserService.On("VerifyEmail", verificationCode).Return(mockResp)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/verifyemail/"+code, nil)

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResp.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})

	t.Run("invalid code", func(t *testing.T) {
		code := "123443"
		verificationCode := utils.Encode(code)
		mockResp := &models.AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusForbidden,
			Message:    "invalid email",
			Err:        errors.New("invalid email"),
		}
		mockUserService.On("VerifyEmail", verificationCode).Return(mockResp)
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/api/auth/verifyemail/"+code, nil)

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResp.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})
}

func TestForgetPasswordController(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCredential := &models.ForgetPasswordInput{
			Email: "bochuang@gmail.com",
		}
		mockResp := &models.AuthServiceResponse{
			Message:    "You will receive a reset email if user with that email exist",
			Status:     "success",
			StatusCode: http.StatusOK,
		}

		reqBody, err := json.Marshal(gin.H{
			"email": mockCredential.Email,
		})

		assert.NoError(t, err)
		mockUserService.On("ForgetPassword", mockCredential.Email).Return(mockResp)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBuffer(reqBody))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResp.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})

	t.Run("no email in request", func(t *testing.T) {

		mockResp := &models.AuthServiceResponse{
			Message:    "invalid request",
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
		}

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/forgotpassword", nil)

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResp.Status,
			"message": mockResp.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResp.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockAuthService.AssertExpectations(t)
	})
}

func TestResetPasswordController(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userCredential := &models.ResetPasswordInput{
			Password:        "12345679",
			PasswordConfirm: "12345679",
		}
		resetToken := "123456"
		mockResponse := &models.AuthServiceResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "Password updated successfully. Please Login with new password",
		}
		mockUserService.On("ResetPassword", userCredential, resetToken).Return(mockResponse)

		reqBody, err := json.Marshal(gin.H{
			"password":        userCredential.Password,
			"passwordConfirm": userCredential.PasswordConfirm,
		})

		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/api/auth/resetpassword/"+resetToken, bytes.NewBuffer(reqBody))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResponse.Status,
			"message": mockResponse.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResponse.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockUserService.AssertExpectations(t)
	})

	t.Run("password not match", func(t *testing.T) {
		userCredential := &models.ResetPasswordInput{
			Password:        "12345679",
			PasswordConfirm: "12345670",
		}
		resetToken := "123456"
		mockResponse := &models.AuthServiceResponse{
			Status:     "fail",
			StatusCode: http.StatusBadRequest,
			Message:    "Password does not match",
		}
		mockUserService.On("ResetPassword", userCredential, resetToken).Return(mockResponse)

		reqBody, err := json.Marshal(gin.H{
			"password":        userCredential.Password,
			"passwordConfirm": userCredential.PasswordConfirm,
		})

		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPatch, "/api/auth/resetpassword/"+resetToken, bytes.NewBuffer(reqBody))

		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		server.ServeHTTP(w, req)

		respBody, err := json.Marshal(gin.H{
			"status":  mockResponse.Status,
			"message": mockResponse.Message,
		})

		assert.NoError(t, err)

		assert.Equal(t, mockResponse.StatusCode, w.Code)
		assert.Equal(t, respBody, w.Body.Bytes())

		mockUserService.AssertExpectations(t)
	})
}
