package services

import (
	"github.com/tonybobo/auth-template/models"
)

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)

	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	ForgetPassword(email string) *AuthServiceResponse
	RefreshAccessToken(cookie string) *AuthServiceResponse
	ResetPassword(user *models.ResetPasswordInput, resetToken string) *AuthServiceResponse
	VerifyEmail(verificationCode string) *AuthServiceResponse
}
