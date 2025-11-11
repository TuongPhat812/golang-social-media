package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myself/golang-social-media/gateway/internal/application/messages"
	"github.com/myself/golang-social-media/gateway/internal/application/users"
)

type createMessageRequest struct {
	SenderID   string `json:"senderId" binding:"required"`
	ReceiverID string `json:"receiverId" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func RegisterRoutes(router *gin.Engine, userService users.Service, messageService messages.Service) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "gateway OK"})
	})

	router.GET("/sample-user", func(c *gin.Context) {
		c.JSON(http.StatusOK, userService.SampleUser())
	})

	router.POST("/chat/messages", func(c *gin.Context) {
		var req createMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		msg, err := messageService.CreateMessage(c.Request.Context(), req.SenderID, req.ReceiverID, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, msg)
	})
}
