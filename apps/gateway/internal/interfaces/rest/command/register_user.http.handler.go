package command

import (
	"net/http"

	app "golang-social-media/apps/gateway/internal/application/command/contracts"
	appdto "golang-social-media/apps/gateway/internal/application/command/dto"
	httpcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

type registerUserHTTPHandler struct {
	command app.RegisterUserCommand
}

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func NewRegisterUserHTTPHandler(command app.RegisterUserCommand) httpcontracts.RegisterUserHTTPHandler {
	return &registerUserHTTPHandler{command: command}
}

func (h *registerUserHTTPHandler) Mount(router *gin.RouterGroup) {
	router.POST("/register", h.handle)
}

func (h *registerUserHTTPHandler) handle(c *gin.Context) {
	var req registerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.command.Handle(c.Request.Context(), appdto.RegisterUserCommandRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		logger.Component("gateway.http.register_user").
			Error().
			Err(err).
			Str("email", req.Email).
			Msg("register user failed")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
