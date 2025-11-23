package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/chat-service/internal/infrastructure/bootstrap"
	grpcserver "golang-social-media/apps/chat-service/internal/infrastructure/grpc"
	interfaces "golang-social-media/apps/chat-service/internal/interfaces/grpc"
	chatgrpc "golang-social-media/apps/chat-service/internal/interfaces/grpc/chat"
	grpcmappers "golang-social-media/apps/chat-service/internal/interfaces/grpc/mappers"
	"golang-social-media/pkg/config"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger.SetModule("chat-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	logger.Component("chat.bootstrap").
		Info().
		Msg("chat service ready")

	// Start event subscribers in background (non-blocking)
	startSubscribers(ctx, deps)

	// Start gRPC server
	port := config.GetEnvInt("CHAT_SERVICE_PORT", 9000)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		// Setup DTO mapper
		messageDTOMapper := grpcmappers.NewMessageDTOMapper()
		handler := chatgrpc.NewHandler(deps, messageDTOMapper)
		chatv1.RegisterChatServiceServer(server, handler)
		interfaces.RegisterServices(server, deps)
	}); err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to serve chat gRPC")
		os.Exit(1)
	}
}

// startSubscribers starts all event subscribers
func startSubscribers(ctx context.Context, deps *bootstrap.Dependencies) {
	if deps.UserSubscriber != nil {
		go deps.UserSubscriber.Consume(ctx)
	}
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.Publisher != nil {
		if err := deps.Publisher.Close(); err != nil {
			logger.Component("chat.bootstrap").
				Error().
				Err(err).
				Msg("failed to close kafka publisher")
		}
	}
	if deps.UserSubscriber != nil {
		if err := deps.UserSubscriber.Close(); err != nil {
			logger.Component("chat.bootstrap").
				Error().
				Err(err).
				Msg("failed to close user subscriber")
		}
	}
	if deps.Cache != nil {
		if err := deps.Cache.Close(); err != nil {
			logger.Component("chat.bootstrap").
				Error().
				Err(err).
				Msg("failed to close cache")
		}
	}
}
