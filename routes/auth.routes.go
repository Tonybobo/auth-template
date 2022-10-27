package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tonybobo/auth-template/controllers"
	"github.com/tonybobo/auth-template/middleware"
	"github.com/tonybobo/auth-template/services"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(userService), rc.authController.LogoutUser)
	router.GET("/verifyemail/:verificationCode", rc.authController.VerifyEmail)
	router.POST("/forgotpassword", rc.authController.ForgetPassword)
	router.PATCH("/resetpassword/:resetToken", rc.authController.ResetPassword)
}
