package http

import (
	"github.com/gin-gonic/gin"
	"github.com/myself/golang-social-media/gateway/internal/application/messages"
	"github.com/myself/golang-social-media/gateway/internal/application/users"
	"github.com/myself/golang-social-media/gateway/internal/interfaces/rest"
)

func NewRouter(userService users.Service, messageService messages.Service) *gin.Engine {
	router := gin.Default()

	rest.RegisterRoutes(router, userService, messageService)

	return router
}
