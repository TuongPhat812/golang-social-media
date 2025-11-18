package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/notification-service/internal/infrastructure/bootstrap"
	grpcserver "golang-social-media/apps/notification-service/internal/infrastructure/grpc"
	interfaces "golang-social-media/apps/notification-service/internal/interfaces/grpc"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"

	"google.golang.org/grpc"
)

func main() {
	logger.SetModule("notification-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	// Start event subscribers
	startSubscribers(ctx, deps)

	logger.Component("notification.bootstrap").
		Info().
		Msg("notification service ready")

	// Start gRPC server
	port := config.GetEnvInt("NOTIFICATION_SERVICE_PORT", 9100)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		interfaces.RegisterServices(server, deps)
	}); err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to serve notification gRPC")
		os.Exit(1)
	}
}

// startSubscribers starts all event subscribers
func startSubscribers(ctx context.Context, deps *bootstrap.Dependencies) {
	deps.ChatSubscriber.Consume(ctx)
	deps.UserSubscriber.Consume(ctx)
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.Publisher != nil {
		if err := deps.Publisher.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close kafka publisher")
		}
	}

	if deps.Session != nil {
		deps.Session.Close()
	}

	if deps.ChatSubscriber != nil {
		if err := deps.ChatSubscriber.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close chat subscriber")
		}
	}

	if deps.UserSubscriber != nil {
		if err := deps.UserSubscriber.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close user subscriber")
		}
	}
}
