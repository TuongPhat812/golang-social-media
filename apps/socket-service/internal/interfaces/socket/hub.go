package socket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/myself/golang-social-media/pkg/events"
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
			log.Printf("failed to upgrade websocket: %v", err)
			return
		}
		defer conn.Close()

		log.Println("socket connected")
		for {
			if _, _, err := conn.NextReader(); err != nil {
				log.Printf("socket closed: %v", err)
				break
			}
		}
	})
}

func (h *Hub) BroadcastChatCreated(event events.ChatCreated) {
	log.Printf("[socket-service] broadcast ChatCreated: %+v", event)
	// TODO: push to connected clients
}

func (h *Hub) BroadcastNotificationCreated(event events.NotificationCreated) {
	log.Printf("[socket-service] broadcast NotificationCreated: %+v", event)
	// TODO: push to connected clients
}
