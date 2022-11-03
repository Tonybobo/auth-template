package services

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/thanhpk/randstr"
	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	AuthRepository models.AuthRepository
	ctx            context.Context
	temp           *template.Template
}

func NewUserServiceImpl(AuthRepository models.AuthRepository, ctx context.Context, temp *template.Template) UserService {
	return &UserServiceImpl{AuthRepository, ctx, temp}
}

func (us *UserServiceImpl) RefreshAccessToken(cookie string) *AuthServiceResponse {
	var result *AuthServiceResponse
	result.Status = "success"
	result.StatusCode = http.StatusOK

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)

	if err != nil {
		result.Status = "fail"
		result.StatusCode = http.StatusForbidden
		result.Message = err.Error()
		result.Err = err
		return result
	}

	oid, err := primitive.ObjectIDFromHex(fmt.Sprint(sub))

	user, err := us.AuthRepository.FindUserById(us.ctx, oid)

	if err != nil {
		result.Err = err
		result.Message = err.Error()
		result.Status = "fail"
		result.StatusCode = http.StatusForbidden
		return result
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)

	if err != nil {
		result.Err = err
		result.Message = err.Error()
		result.Status = "fail"
		result.StatusCode = http.StatusForbidden
		return result
	}

	result.AccessToken = access_token
	return result

}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid User ID")
	}

	user, err := us.AuthRepository.FindUserById(us.ctx, oid)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	user, err := us.AuthRepository.FindUserByEmail(us.ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) UpdateUserById(id, field, value string) (*models.DBResponse, error) {
	userId, _ := primitive.ObjectIDFromHex(id)

	_, err := us.AuthRepository.UpdateUserById(us.ctx, userId, field, value)

	if err != nil {
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}

func (us *UserServiceImpl) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {

	_, err := us.AuthRepository.UpdateOne(us.ctx, field, value)

	if err != nil {
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}

func (us *UserServiceImpl) VerifyEmail(verificationCode string) (*mongo.UpdateResult, error) {

	result, err := us.AuthRepository.VerifyEmail(us.ctx, verificationCode)

	return result, err
}

func (us *UserServiceImpl) ClearResetPasswordToken(token string, password string) (*mongo.UpdateResult, error) {

	result, err := us.AuthRepository.ClearResetPasswordToken(us.ctx, token, password)

	return result, err
}

func (us *UserServiceImpl) ForgetPassword(email string) *AuthServiceResponse {
	var response *AuthServiceResponse

	response.Message = "You will receive a reset email if user with that email exist"
	response.Status = "success"

	user, err := us.AuthRepository.FindUserByEmail(us.ctx, email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.StatusCode = http.StatusOK
		}
		response.StatusCode = http.StatusBadGateway
		response.Status = "fail"
		response.Message = err.Error()
		response.Err = err
		return response
	}

	if !user.Verified {
		response.StatusCode = http.StatusUnauthorized
		response.Status = "fail"
		return response
	}

	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	result, err := us.AuthRepository.ResetPasswordToken(us.ctx, email, passwordResetToken)

	if err != nil {
		response.StatusCode = http.StatusForbidden
		response.Status = "fail"
		response.Message = "There was a error sending reset email"
		response.Err = err
		return response
	}

	if result.MatchedCount == 0 {
		response.StatusCode = http.StatusOK
		response.Status = "success"
		return response
	}

	var firstName = user.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	config, _ := config.LoadConfig(".")

	emailData := utils.EmailData{
		URL:       "http://localhost:" + config.Port + "/api/auth/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Please Reset the password within 15 minutes",
	}

	err = utils.SendEmail(user, &emailData, us.temp, "resetPassword.html")

	if err != nil {
		response.StatusCode = http.StatusBadGateway
		response.Err = err
		response.Status = "fail"
		response.Message = err.Error()
		return response
	}

	return response

}
