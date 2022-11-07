package services

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServiceImpl struct {
	AuthRepository models.AuthRepository
	ctx            context.Context
	temp           *template.Template
}

type AuthServiceResponse struct {
	User               *models.DBResponse
	Status             string
	Err                error
	Message            string
	StatusCode         int
	AccessToken        string
	RefreshAccessToken string
}

func NewAuthService(AuthRepository models.AuthRepository, ctx context.Context, temp *template.Template) AuthService {
	return &AuthServiceImpl{AuthRepository, ctx, temp}
}

func (uc *AuthServiceImpl) Test() *AuthServiceResponse {

	return &AuthServiceResponse{
		Message: "test",
	}
}

func (uc *AuthServiceImpl) SignInUser(credential *models.SignInInput) *AuthServiceResponse {

	result := &AuthServiceResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
	}

	user, err := uc.AuthRepository.FindUserByEmail(uc.ctx, credential.Email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			result.Status = "fail"
			result.StatusCode = http.StatusBadRequest
			result.Message = "Invalid Email or Password"
			result.Err = err
			return result
		}
		result.Status = "fail"
		result.StatusCode = http.StatusBadGateway
		result.Message = err.Error()
		result.Err = err
		return result
	}

	if !user.Verified {
		result.Err = errors.New("you have not verify the account , please verify your email to login ")
		result.Status = "fail"
		result.StatusCode = http.StatusUnauthorized
		result.Message = result.Err.Error()
		return result
	}

	if err := utils.VerifyPassword(user.Password, credential.Password); err != nil {
		result.Status = "fail"
		result.StatusCode = http.StatusUnauthorized
		result.Message = "Invalid Email or Password"
		result.Err = err
		return result
	}

	config, _ := config.LoadConfig(".")

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)

	if err != nil {
		result.Status = "fail"
		result.StatusCode = http.StatusBadRequest
		result.Message = err.Error()
		result.Err = err
		return result
	}

	result.AccessToken = access_token

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)

	if err != nil {
		result.Status = "fail"
		result.StatusCode = http.StatusBadRequest
		result.Err = err
		result.Message = err.Error()
		return result
	}

	result.RefreshAccessToken = refresh_token

	return result
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) *AuthServiceResponse {

	result := &AuthServiceResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
	}

	if user.Password != user.PasswordConfirm {
		result.Err = errors.New("password not match")
		result.Message = "password not match"
		result.Status = "fail"
		result.StatusCode = http.StatusBadRequest
		return result
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	newUser, code, err := uc.AuthRepository.SignUpUser(uc.ctx, user)
	if err != nil {
		result.Err = err
		result.Status = "fail"
		result.Message = err.Error()
		result.StatusCode = http.StatusBadGateway
		return result
	}

	result.User = newUser
	result.Message = "An email with the verification code has been sent to " + result.User.Email

	if err != nil {
		result.Err = err
		result.Status = "fail"
		result.Message = err.Error()
		result.StatusCode = http.StatusBadGateway
		return result
	}

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	config, err := config.LoadConfig("../")

	if err != nil {
		log.Panic("could not load environment variables")
	}

	emailData := utils.EmailData{
		URL:       "http://localhost:" + config.Port + "/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Please Verify",
	}

	err = utils.SendEmail(newUser, &emailData, uc.temp, "verification.html")

	if err != nil {
		result.Err = err
		result.Message = err.Error()
		return result
	}

	return result
}
