package command

import (
	"net/http"

	app "golang-social-media/apps/gateway/internal/application/command/contracts"
	httpcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

type createMessageHTTPHandler struct {
	command app.CreateMessageCommand
}

type createMessageRequest struct {
	SenderID   string `json:"senderId" binding:"required"`
	ReceiverID string `json:"receiverId" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func NewCreateMessageHTTPHandler(command app.CreateMessageCommand) httpcontracts.CreateMessageHTTPHandler {
	return &createMessageHTTPHandler{command: command}
}

func (h *createMessageHTTPHandler) Mount(router *gin.RouterGroup) {
	router.POST("/chat/messages", h.handle)
}

func (h *createMessageHTTPHandler) handle(c *gin.Context) {
	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.command.Handle(c.Request.Context(), req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		logger.Component("gateway.http.create_message").
			Error().
			Err(err).
			Msg("create message failed")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}
