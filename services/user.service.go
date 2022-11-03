package services

import (
	"github.com/tonybobo/auth-template/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id, field, value string) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	ForgetPassword(email string) *AuthServiceResponse
	RefreshAccessToken(cookie string) *AuthServiceResponse
	ClearResetPasswordToken(resetPasswordToken string, hashPasswordToken string) (*mongo.UpdateResult, error)
	VerifyEmail(verificationCode string) (*mongo.UpdateResult, error)
}
