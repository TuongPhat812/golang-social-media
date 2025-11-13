package rest

import (
	"net/http"

	command "golang-social-media/apps/auth-service/internal/application/command"
	query "golang-social-media/apps/auth-service/internal/application/query"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	RegisterUser *command.RegisterUserHandler
	LoginUser    *command.LoginUserHandler
	GetProfile   *query.GetUserProfileHandler
}

func NewRouter(h Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.POST("/auth/register", func(c *gin.Context) {
		var req auth.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, auth.ErrorResponse{Error: err.Error()})
			return
		}
		resp, err := h.RegisterUser.Handle(c.Request.Context(), req)
		if err != nil {
			logger.Component("auth.register").
				Error().
				Err(err).
				Str("email", req.Email).
				Msg("register failed")
			c.JSON(http.StatusBadRequest, auth.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	})

	router.POST("/auth/login", func(c *gin.Context) {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, auth.ErrorResponse{Error: err.Error()})
			return
		}
		resp, err := h.LoginUser.Handle(c.Request.Context(), req)
		if err != nil {
			logger.Component("auth.login").
				Error().
				Err(err).
				Str("email", req.Email).
				Msg("login failed")
			c.JSON(http.StatusUnauthorized, auth.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.GET("/auth/profile/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := h.GetProfile.Handle(c.Request.Context(), id)
		if err != nil {
			logger.Component("auth.profile").
				Error().
				Err(err).
				Str("user_id", id).
				Msg("profile lookup failed")
			c.JSON(http.StatusNotFound, auth.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "auth-service OK"})
	})

	return router
}
