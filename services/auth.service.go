package services

import "github.com/tonybobo/auth-template/models"

type AuthService interface {
	Test() *models.AuthServiceResponse
	SignUpUser(user *models.SignUpInput) *models.AuthServiceResponse
	SignInUser(user *models.SignInInput) *models.AuthServiceResponse
}
