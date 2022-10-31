package controllers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/services"
	"github.com/tonybobo/auth-template/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
	temp        *template.Template
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection, temp *template.Template) AuthController {
	return AuthController{authService, userService, ctx, collection, temp}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.BindJSON(&user); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	code := randstr.String(20)
	verificationCode := utils.Encode(code)

	ac.userService.UpdateUserById(newUser.ID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.EmailData{
		URL:       "http://localhost:8080/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Please Verify",
	}

	err = utils.SendEmail(newUser, &emailData, ac.temp, "verification.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "An email with the verification code has been sent to " + newUser.Email
	ctx.JSON(http.StatusOK, gin.H{"status": "Success", "message": message})

}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credential *models.SignInInput

	if err := ctx.ShouldBindJSON(&credential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByEmail(credential.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Email or Password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You have not verify the account , Please verify your email to login "})
		return
	}

	if err := utils.VerifyPassword(user.Password, credential.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid Email or Password"})
		return
	}

	config, _ := config.LoadConfig(".")

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token, "refresh_token": refresh_token})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "The user with this id no longer exist"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	fmt.Print(ctx.Get("currentUser"))
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}},
		{Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}

	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Invalid Email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully Verified"})
}

func (ac *AuthController) ForgetPassword(ctx *gin.Context) {
	var credential models.ForgetPasswordInput
	if err := ctx.ShouldBindJSON(&credential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(credential.Email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusOK, gin.H{"status": "Success", "message": message})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail ", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Account not verified"})
		return
	}

	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	query := bson.D{{Key: "email", Value: strings.ToLower(user.Email)}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "passwordResetToken", Value: passwordResetToken},
			{Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)},
		}},
	}

	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "There was a error sending reset email"})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Reset Email has been sent"})
		return
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

	err = utils.SendEmail(user, &emailData, ac.temp, "resetPassword.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "There was an error sending email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var userCredential *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userCredential.Password != userCredential.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password does not match"})
		return
	}

	hashPassword, _ := utils.HashPassword(userCredential.Password)
	resetPasswordToken := utils.Encode(resetToken)

	query := bson.D{{Key: "passwordResetToken", Value: resetPasswordToken}, {Key: "passwordResetAt", Value: bson.D{{Key: "$gt", Value: time.Now()}}}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "password", Value: hashPassword}}},
		{Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}

	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Invalid or Expired Token"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password updated successfully. Please Login with new password"})
}
