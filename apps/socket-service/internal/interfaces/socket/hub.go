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
			logger.Error().Err(err).Msg("socket-service failed to upgrade websocket")
			return
		}
		defer conn.Close()

		logger.Info().Msg("socket connected")
		for {
			if _, _, err := conn.NextReader(); err != nil {
				logger.Info().Err(err).Msg("socket connection closed")
				break
			}
		}
	})
}

func (h *Hub) BroadcastChatCreated(event events.ChatCreated) {
	logger.Info().
		Str("topic", events.TopicChatCreated).
		Str("message_id", event.Message.ID).
		Msg("socket broadcast chat update")
	// TODO: push to connected clients
}

func (h *Hub) BroadcastNotificationCreated(event events.NotificationCreated) {
	logger.Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.Notification.ID).
		Str("user_id", event.Notification.UserID).
		Msg("socket broadcast notification update")
	// TODO: push to connected clients
}
