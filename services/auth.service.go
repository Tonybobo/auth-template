package services

import "github.com/tonybobo/auth-template/models"

type AuthService interface {
	Test() *AuthServiceResponse
	SignUpUser(user *models.SignUpInput) *AuthServiceResponse
	SignInUser(user *models.SignInInput) *AuthServiceResponse
}
