package services

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

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

func (us *UserServiceImpl) RefreshAccessToken(cookie string) *models.AuthServiceResponse {

	result := &models.AuthServiceResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
	}

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

	if err != nil {
		result.Err = err
		result.Message = "Invalid Token"
		result.Status = "fail"
		result.StatusCode = http.StatusForbidden
		return result
	}

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

func (us *UserServiceImpl) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {

	_, err := us.AuthRepository.UpdateOne(us.ctx, field, value)

	if err != nil {
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}

func (us *UserServiceImpl) VerifyEmail(verificationCode string) *models.AuthServiceResponse {

	response := &models.AuthServiceResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Message:    "Successfully Verified",
	}
	fmt.Print(verificationCode)
	err := us.AuthRepository.VerifyEmail(us.ctx, verificationCode)

	if err != nil {
		if err.Error() == "invalid email" {
			response.Err = err
			response.Message = err.Error()
			response.Status = "fail"
			response.StatusCode = http.StatusForbidden
			return response
		}
		response.Err = err
		response.Message = err.Error()
		response.Status = "fail"
		response.StatusCode = http.StatusBadGateway
		return response
	}

	return response
}

func (us *UserServiceImpl) ResetPassword(user *models.ResetPasswordInput, resetToken string) *models.AuthServiceResponse {

	response := &models.AuthServiceResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Message:    "Password updated successfully. Please Login with new password",
	}

	if user.Password != user.PasswordConfirm {
		response.Status = "fail"
		response.StatusCode = http.StatusBadRequest
		response.Message = "Password does not match"

		return response
	}

	err := us.AuthRepository.ClearResetPasswordToken(us.ctx, resetToken, user.Password)

	if err != nil {
		if err.Error() == "invalid or expired token" {
			response.Status = "fail"
			response.StatusCode = http.StatusForbidden
			response.Message = err.Error()
			response.Err = err
			return response
		}
		response.Status = "fail"
		response.StatusCode = http.StatusBadGateway
		response.Message = err.Error()
		response.Err = err

		return response
	}

	return response
}

func (us *UserServiceImpl) ForgetPassword(email string) *models.AuthServiceResponse {

	response := &models.AuthServiceResponse{
		Message:    "You will receive a reset email if user with that email exist",
		Status:     "success",
		StatusCode: http.StatusOK,
	}

	user, resetToken, err := us.AuthRepository.ForgetPassword(us.ctx, email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			response.StatusCode = http.StatusOK
			return response
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
		response.Err = errors.New("account has not been verified. please verify your account with the email sent")
		response.Message = response.Err.Error()
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
