package command

import (
	"net/http"
	"time"

	app "golang-social-media/apps/gateway/internal/application/command/contracts"
	"golang-social-media/apps/gateway/internal/infrastructure/middleware"
	httpcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/command/contracts"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type createMessageHTTPHandler struct {
	command app.CreateMessageCommand
	log     *zerolog.Logger
}

type createMessageRequest struct {
	ReceiverID string `json:"receiverId" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func NewCreateMessageHTTPHandler(command app.CreateMessageCommand) httpcontracts.CreateMessageHTTPHandler {
	return &createMessageHTTPHandler{
		command: command,
		log:     logger.Component("gateway.http.create_message"),
	}
}

func (h *createMessageHTTPHandler) Mount(router *gin.RouterGroup) {
	router.POST("/chat/messages", h.handle)
}

func (h *createMessageHTTPHandler) handle(c *gin.Context) {
	startTime := time.Now()

	// Get user ID from JWT middleware
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.log.Warn().
			Msg("user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	senderID, ok := userID.(string)
	if !ok || senderID == "" {
		h.log.Warn().
			Msg("invalid user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}

	// Parse request
	parseStart := time.Now()
	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		parseDuration := time.Since(parseStart)
		totalDuration := time.Since(startTime)
		h.log.Warn().
			Err(err).
			Dur("parse_request_ms", parseDuration).
			Dur("total_ms", totalDuration).
			Msg("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parseDuration := time.Since(parseStart)

	// Execute command (sender_id comes from JWT, not request body)
	commandStart := time.Now()
	msg, err := h.command.Handle(c.Request.Context(), senderID, req.ReceiverID, req.Content)
	commandDuration := time.Since(commandStart)

	if err != nil {
		totalDuration := time.Since(startTime)
		h.log.Error().
			Err(err).
			Str("sender_id", senderID).
			Str("receiver_id", req.ReceiverID).
			Dur("parse_request_ms", parseDuration).
			Dur("command_exec_ms", commandDuration).
			Dur("total_ms", totalDuration).
			Msg("create message failed")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Serialize response
	serializeStart := time.Now()
	c.JSON(http.StatusCreated, msg)
	serializeDuration := time.Since(serializeStart)

	totalDuration := time.Since(startTime)

	h.log.Info().
		Str("message_id", msg.ID).
		Dur("parse_request_ms", parseDuration).
		Dur("command_exec_ms", commandDuration).
		Dur("serialize_response_ms", serializeDuration).
		Dur("total_ms", totalDuration).
		Msg("HTTP request completed")
}
