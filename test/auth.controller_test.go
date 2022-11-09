package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tonybobo/auth-template/controllers"
	"github.com/tonybobo/auth-template/mocks"
	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/routes"
)

var (
	mockAuthService     = new(mocks.MockAuthService)
	mockUserService     = new(mocks.MockUserService)
	ctx                 = context.TODO()
	authController      = controllers.NewAuthController(mockAuthService, mockUserService, ctx)
	authRouteController = routes.NewAuthRouteController(authController)
	server              = gin.Default()
	router              = server.Group("/api")
)

func TestAuth(t *testing.T) {
	mockAuthService.On("Test").Return(&models.AuthServiceResponse{
		Message: "test",
	})
	authRouteController.AuthRoute(router, mockUserService)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/auth/test", nil)

	if err != nil {
		t.FailNow()
	}
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":0}", w.Body.String())
}
