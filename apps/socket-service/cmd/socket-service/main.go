package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	appevents "golang-social-media/apps/socket-service/internal/application/events"
	"golang-social-media/apps/socket-service/internal/infrastructure/eventbus"
	"golang-social-media/apps/socket-service/internal/interfaces/socket"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
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
		logger.Error().Err(err).Msg("failed to create socket kafka listener")
		os.Exit(1)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close socket kafka listener")
		}
	}()

	listener.Start(ctx)

	router := gin.Default()
	hub.RegisterRoutes(router)

	port := config.GetEnvInt("SOCKET_SERVICE_PORT", 9200)
	addr := fmt.Sprintf(":%d", port)

	logger.Info().Str("addr", addr).Msg("socket service starting")
	if err := router.Run(addr); err != nil {
		logger.Error().Err(err).Msg("failed to start socket service")
		os.Exit(1)
	}
}
