package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/auth-service/internal/infrastructure/bootstrap"
	"golang-social-media/apps/auth-service/internal/interfaces/rest"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
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

	// Setup router
	router := rest.NewRouter(rest.Handlers{
		RegisterUser: deps.RegisterUserCmd,
		LoginUser:    deps.LoginUserCmd,
		GetProfile:   deps.GetUserProfileQuery,
	})

	port := config.GetEnvInt("AUTH_SERVICE_PORT", 9101)
	addr := fmt.Sprintf(":%d", port)

	logger.Component("auth.http").
		Info().
		Str("addr", addr).
		Msg("auth service starting")

	if err := router.Run(addr); err != nil {
		logger.Component("auth.http").
			Error().
			Err(err).
			Msg("auth service failed")
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
}
