package http

import (
	"net/http"

	commandcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	querycontracts "golang-social-media/apps/gateway/internal/interfaces/rest/query/contracts"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	registerUser commandcontracts.RegisterUserHTTPHandler,
	loginUser commandcontracts.LoginUserHTTPHandler,
	createMessage commandcontracts.CreateMessageHTTPHandler,
	getUserProfile querycontracts.GetUserProfileHTTPHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "gateway OK"})
	})

	authGroup := router.Group("/auth")
	registerUser.Mount(authGroup)
	loginUser.Mount(authGroup)
	getUserProfile.Mount(authGroup)

	apiGroup := router.Group("")
	createMessage.Mount(apiGroup)

	return router
}
