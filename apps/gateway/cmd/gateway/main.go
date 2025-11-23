package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/gateway/internal/infrastructure/bootstrap"
	httpserver "golang-social-media/apps/gateway/internal/infrastructure/http"
	middleware "golang-social-media/apps/gateway/internal/infrastructure/middleware"
	commandrest "golang-social-media/apps/gateway/internal/interfaces/rest/command"
	queryrest "golang-social-media/apps/gateway/internal/interfaces/rest/query"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.SetModule("gateway")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("gateway.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	logger.Component("gateway.bootstrap").
		Info().
		Msg("gateway service ready")

	// Setup router
	router := buildRouter(deps)

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

func buildRouter(deps *bootstrap.Dependencies) *gin.Engine {
	createMessageHTTP := commandrest.NewCreateMessageHTTPHandler(deps.CreateMessageCmd)
	registerUserHTTP := commandrest.NewRegisterUserHTTPHandler(deps.RegisterUserCmd)
	loginUserHTTP := commandrest.NewLoginUserHTTPHandler(deps.LoginUserCmd)
	getUserProfileHTTP := queryrest.NewGetUserProfileHTTPHandler(deps.GetUserProfileQuery)

	// Create auth client adapter for middleware
	authClientAdapter := middleware.NewAuthGRPCClientAdapter(deps.AuthGRPCClient)

	return httpserver.NewRouter(
		registerUserHTTP,
		loginUserHTTP,
		createMessageHTTP,
		getUserProfileHTTP,
		authClientAdapter,
	)
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.ChatClient != nil {
		if err := deps.ChatClient.Close(); err != nil {
			logger.Component("gateway.bootstrap").
				Error().
				Err(err).
				Msg("failed to close chat client")
		}
	}
	if deps.AuthGRPCClient != nil {
		if err := deps.AuthGRPCClient.Close(); err != nil {
			logger.Component("gateway.bootstrap").
				Error().
				Err(err).
				Msg("failed to close auth gRPC client")
		}
	}
}
