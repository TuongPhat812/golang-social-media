package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/auth-service/internal/infrastructure/bootstrap"
	grpcserver "golang-social-media/apps/auth-service/internal/infrastructure/grpc"
	"golang-social-media/apps/auth-service/internal/interfaces/grpc"
	authgrpc "golang-social-media/apps/auth-service/internal/interfaces/grpc/auth"
	"golang-social-media/apps/auth-service/internal/interfaces/rest"
	authv1 "golang-social-media/pkg/gen/auth/v1"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger.SetModule("auth-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("auth.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	logger.Component("auth.bootstrap").
		Info().
		Msg("auth service ready")

	// Setup HTTP handlers
	authHandler := rest.NewAuthHandler(deps.RegisterUserCmd, deps.LoginUserCmd)
	profileHandler := rest.NewProfileHandler(deps.UpdateProfileCmd, deps.GetUserProfileQuery, deps.GetCurrentUserQuery)
	passwordHandler := rest.NewPasswordHandler(deps.ChangePasswordCmd)
	tokenHandler := rest.NewTokenHandler(deps.LogoutUserCmd, deps.RefreshTokenCmd, deps.RevokeTokenCmd, deps.ValidateTokenQuery)

	handlers := rest.NewHandlers(authHandler, profileHandler, passwordHandler, tokenHandler)

	// Setup HTTP router
	router := rest.NewRouter(handlers, deps.JwtService, deps.Cache)

	// Start HTTP server in goroutine
	httpPort := config.GetEnvInt("AUTH_SERVICE_PORT", 9101)
	httpAddr := fmt.Sprintf(":%d", httpPort)

	go func() {
		logger.Component("auth.http").
			Info().
			Str("addr", httpAddr).
			Msg("auth HTTP server starting")

		if err := router.Run(httpAddr); err != nil && err != http.ErrServerClosed {
			logger.Component("auth.http").
				Error().
				Err(err).
				Msg("auth HTTP server failed")
			os.Exit(1)
		}
	}()

	// Start gRPC server
	grpcPort := config.GetEnvInt("AUTH_SERVICE_GRPC_PORT", 9100)
	grpcAddr := fmt.Sprintf(":%d", grpcPort)

	if err := grpcserver.Start(grpcAddr, func(server *grpc.Server) {
		// Setup gRPC handler
		handler := authgrpc.NewHandler(deps)
		authv1.RegisterAuthServiceServer(server, handler)
		grpc.RegisterServices(server, deps)
	}); err != nil {
		logger.Component("auth.grpc").
			Error().
			Err(err).
			Msg("auth gRPC server failed")
		os.Exit(1)
	}
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.Publisher != nil {
		if err := deps.Publisher.Close(); err != nil {
			logger.Component("auth.bootstrap").
				Error().
				Err(err).
				Msg("failed to close kafka publisher")
		}
	}
	if deps.Cache != nil {
		if err := deps.Cache.Close(); err != nil {
			logger.Component("auth.bootstrap").
				Error().
				Err(err).
				Msg("failed to close cache")
		}
	}
}
