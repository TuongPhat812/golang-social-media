package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang-social-media/apps/gateway/internal/application/messages"
	"golang-social-media/apps/gateway/internal/application/users"
	chatclient "golang-social-media/apps/gateway/internal/infrastructure/grpc/chat"
	httpserver "golang-social-media/apps/gateway/internal/infrastructure/http"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.SetModule("gateway")
	config.LoadEnv()

	switch strings.ToLower(config.GetEnv("GIN_MODE", "debug")) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if strings.EqualFold(config.GetEnv("GIN_DISABLE_ACCESS_LOG", "false"), "true") {
		gin.DefaultWriter = io.Discard
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	chatClient, err := chatclient.NewClient(ctx)
	if err != nil {
		logger.Component("gateway.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect to chat service")
		os.Exit(1)
	}
	defer func() {
		if err := chatClient.Close(); err != nil {
			logger.Component("gateway.grpc").
				Error().
				Err(err).
				Msg("failed to close chat client")
		}
	}()

	userService := users.NewService()
	messageService := messages.NewService(chatClient, logger.Component("gateway.messages"))

	router := httpserver.NewRouter(userService, messageService)

	port := config.GetEnvInt("GATEWAY_PORT", 8080)
	addr := fmt.Sprintf(":%d", port)

	logger.Component("gateway.http").
		Info().
		Str("addr", addr).
		Msg("gateway service starting")
	if err := router.Run(addr); err != nil {
		logger.Component("gateway.http").
			Error().
			Err(err).
			Msg("failed to start gateway")
		os.Exit(1)
	}
}
