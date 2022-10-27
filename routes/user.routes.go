package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tonybobo/auth-template/controllers"
	"github.com/tonybobo/auth-template/middleware"
	"github.com/tonybobo/auth-template/services"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", uc.userController.GetMe)
}
