package command

import (
	"net/http"

	app "golang-social-media/apps/gateway/internal/application/command/contracts"
	appdto "golang-social-media/apps/gateway/internal/application/command/dto"
	httpcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

type loginUserHTTPHandler struct {
	command app.LoginUserCommand
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewLoginUserHTTPHandler(command app.LoginUserCommand) httpcontracts.LoginUserHTTPHandler {
	return &loginUserHTTPHandler{command: command}
}

func (h *loginUserHTTPHandler) Mount(router *gin.RouterGroup) {
	router.POST("/login", h.handle)
}

func (h *loginUserHTTPHandler) handle(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.command.Handle(c.Request.Context(), appdto.LoginUserCommandRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Component("gateway.http.login_user").
			Error().
			Err(err).
			Str("email", req.Email).
			Msg("login user failed")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": resp.UserID,
		"token":  resp.Token,
	})
}
