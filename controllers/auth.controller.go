package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/services"
	"github.com/tonybobo/auth-template/utils"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context) AuthController {
	return AuthController{authService, userService, ctx}
}

func (ac *AuthController) Test(ctx *gin.Context) {
	response := ac.authService.Test()
	ctx.JSON(http.StatusOK, gin.H{"message": response.StatusCode})
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.BindJSON(&user); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	response := ac.authService.SignUpUser(user)

	ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})

}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credential *models.SignInInput

	if err := ctx.ShouldBindJSON(&credential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	config, _ := config.LoadConfig(".")

	response := ac.authService.SignInUser(credential)

	if response.Err != nil {
		ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})
		return
	}

	ctx.SetCookie("access_token", response.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", response.RefreshAccessToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": response.AccessToken, "refresh_token": response.RefreshAccessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	response := ac.userService.RefreshAccessToken(cookie)

	if response.Err != nil {
		ctx.AbortWithStatusJSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})
		return
	}

	ctx.SetCookie("access_token", response.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "access_token": response.AccessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	response := ac.userService.VerifyEmail(verificationCode)

	ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})
}

func (ac *AuthController) ForgetPassword(ctx *gin.Context) {
	var credential models.ForgetPasswordInput
	if err := ctx.ShouldBindJSON(&credential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	response := ac.userService.ForgetPassword(credential.Email)

	ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var userCredential *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	response := ac.userService.ResetPassword(userCredential, resetToken)

	if response.Err != nil {
		ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(response.StatusCode, gin.H{"status": response.Status, "message": response.Message})
}
