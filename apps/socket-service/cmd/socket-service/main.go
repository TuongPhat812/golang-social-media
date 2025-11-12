package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang-social-media/pkg/config"
	appevents "golang-social-media/apps/socket-service/internal/application/events"
	"golang-social-media/apps/socket-service/internal/infrastructure/eventbus"
	"golang-social-media/apps/socket-service/internal/interfaces/socket"
)

func main() {
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	hub := socket.NewHub()
	eventService := appevents.NewService(hub)

	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	listener, err := eventbus.NewListener(
		brokers,
		config.GetEnv("SOCKET_CHAT_GROUP_ID", "socket-service-chat"),
		config.GetEnv("SOCKET_NOTIFICATION_GROUP_ID", "socket-service-notification"),
		eventService,
	)
	if err != nil {
		log.Fatalf("failed to create kafka listener: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("failed to close kafka listener: %v", err)
		}
	}()

	listener.Start(ctx)

	router := gin.Default()
	hub.RegisterRoutes(router)

	port := config.GetEnvInt("SOCKET_SERVICE_PORT", 9200)
	addr := fmt.Sprintf(":%d", port)

	log.Printf("socket service starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start socket service: %v", err)
	}
}
