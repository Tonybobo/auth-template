package services

import (
	"github.com/tonybobo/auth-template/models"
)

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)

	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	ForgetPassword(email string) *models.AuthServiceResponse
	RefreshAccessToken(cookie string) *models.AuthServiceResponse
	ResetPassword(user *models.ResetPasswordInput, resetToken string) *models.AuthServiceResponse
	VerifyEmail(verificationCode string) *models.AuthServiceResponse
}
