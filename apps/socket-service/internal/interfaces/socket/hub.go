package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type Hub struct {
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Hub) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Component("socket.hub").
				Error().
				Err(err).
				Msg("failed to upgrade websocket")
			return
		}
		defer conn.Close()

		logger.Component("socket.hub").
			Info().
			Msg("socket connected")
		for {
			if _, _, err := conn.NextReader(); err != nil {
				logger.Component("socket.hub").
					Info().
					Err(err).
					Msg("socket connection closed")
				break
			}
		}
	})
}

func (h *Hub) BroadcastChatCreated(event events.ChatCreated) {
	logger.Component("socket.hub").
		Info().
		Str("topic", events.TopicChatCreated).
		Str("message_id", event.Message.ID).
		Msg("broadcast chat update")
	// TODO: push to connected clients
}

func (h *Hub) BroadcastNotificationCreated(event events.NotificationCreated) {
	logger.Component("socket.hub").
		Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.Notification.ID).
		Str("user_id", event.Notification.UserID).
		Msg("broadcast notification update")
	// TODO: push to connected clients
}
