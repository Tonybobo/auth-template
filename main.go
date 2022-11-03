package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/controllers"
	"github.com/tonybobo/auth-template/repository"
	"github.com/tonybobo/auth-template/routes"
	"github.com/tonybobo/auth-template/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	temp *template.Template
)

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err = mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Print("Successfully connected to DB")

	authRepository := repository.NewAuthRepository(mongoClient.Database("golang_mongodb").Collection("users"))
	userService = services.NewUserServiceImpl(authRepository, ctx, temp)
	authService = services.NewAuthService(authRepository, ctx, temp)

	AuthController = controllers.NewAuthController(authService, userService, ctx)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewUserRouteController(UserController)

	server = gin.Default()

}

func SetUpRouter() *gin.Engine {

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "server connected"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	return server
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	defer mongoClient.Disconnect(ctx)
	server := SetUpRouter()

	log.Fatal(server.Run(":" + config.Port))
}
