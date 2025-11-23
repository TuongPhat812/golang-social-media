package http

import (
	"net/http"

	"golang-social-media/apps/gateway/internal/infrastructure/middleware"
	commandcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	querycontracts "golang-social-media/apps/gateway/internal/interfaces/rest/query/contracts"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	registerUser commandcontracts.RegisterUserHTTPHandler,
	loginUser commandcontracts.LoginUserHTTPHandler,
	createMessage commandcontracts.CreateMessageHTTPHandler,
	getUserProfile querycontracts.GetUserProfileHTTPHandler,
	authClient middleware.AuthClient,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "gateway OK"})
	})

	// Public routes (no auth required)
	authGroup := router.Group("/auth")
	registerUser.Mount(authGroup)
	loginUser.Mount(authGroup)
	getUserProfile.Mount(authGroup)

	// Protected routes (JWT required - validated via auth service gRPC)
	apiGroup := router.Group("")
	apiGroup.Use(middleware.JWTAuthMiddleware(authClient))
	createMessage.Mount(apiGroup)

	return router
}
