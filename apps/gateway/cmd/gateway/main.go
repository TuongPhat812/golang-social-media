package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang-social-media/apps/gateway/internal/application/messages"
	"golang-social-media/apps/gateway/internal/application/users"
	chatclient "golang-social-media/apps/gateway/internal/infrastructure/grpc/chat"
	httpserver "golang-social-media/apps/gateway/internal/infrastructure/http"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

func main() {
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	chatClient, err := chatclient.NewClient(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to chat service")
		os.Exit(1)
	}
	defer func() {
		if err := chatClient.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close chat client")
		}
	}()

	userService := users.NewService()
	messageService := messages.NewService(chatClient)

	router := httpserver.NewRouter(userService, messageService)

	port := config.GetEnvInt("GATEWAY_PORT", 8080)
	addr := fmt.Sprintf(":%d", port)

	logger.Info().Str("addr", addr).Msg("gateway service starting")
	if err := router.Run(addr); err != nil {
		logger.Error().Err(err).Msg("failed to start gateway")
		os.Exit(1)
	}
}
