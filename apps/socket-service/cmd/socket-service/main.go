package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	bootstrap "golang-social-media/apps/socket-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

func main() {
	logger.SetModule("socket-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("socket.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	logger.Component("socket.bootstrap").
		Info().
		Msg("socket service ready")

	// Start event subscribers in background (non-blocking)
	startSubscribers(ctx, deps)

	// Setup router
	router := gin.Default()
	deps.Hub.RegisterRoutes(router)

	port := config.GetEnvInt("SOCKET_SERVICE_PORT", 9200)
	addr := fmt.Sprintf(":%d", port)

	logger.Component("socket.http").
		Info().
		Str("addr", addr).
		Msg("socket service starting")

	if err := router.Run(addr); err != nil {
		logger.Component("socket.http").
			Error().
			Err(err).
			Msg("failed to start socket service")
		os.Exit(1)
	}
}

// startSubscribers starts all event subscribers
func startSubscribers(ctx context.Context, deps *bootstrap.Dependencies) {
	go deps.ChatSubscriber.Consume(ctx)
	go deps.NotificationSubscriber.Consume(ctx)
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.ChatSubscriber != nil {
		if err := deps.ChatSubscriber.Close(); err != nil {
			logger.Component("socket.bootstrap").
				Error().
				Err(err).
				Msg("failed to close chat subscriber")
		}
	}

	if deps.NotificationSubscriber != nil {
		if err := deps.NotificationSubscriber.Close(); err != nil {
			logger.Component("socket.bootstrap").
				Error().
				Err(err).
				Msg("failed to close notification subscriber")
		}
	}
}
