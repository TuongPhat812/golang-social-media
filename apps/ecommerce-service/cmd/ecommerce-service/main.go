package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bootstrap "golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	grpcserver "golang-social-media/apps/ecommerce-service/internal/infrastructure/grpc"
	interfaces "golang-social-media/apps/ecommerce-service/internal/interfaces/grpc"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger.SetModule("ecommerce-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup all dependencies
	deps, err := bootstrap.SetupDependencies(ctx)
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to setup dependencies")
		os.Exit(1)
	}
	defer cleanup(deps)

	logger.Component("ecommerce.bootstrap").
		Info().
		Msg("ecommerce service ready")

	// Start gRPC server
	port := config.GetEnvInt("ECOMMERCE_SERVICE_PORT", 9200)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		interfaces.RegisterServices(server, deps)
	}); err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to serve ecommerce gRPC")
		os.Exit(1)
	}
}

// cleanup closes all resources
func cleanup(deps *bootstrap.Dependencies) {
	if deps.Publisher != nil {
		if err := deps.Publisher.Close(); err != nil {
			logger.Component("ecommerce.bootstrap").
				Error().
				Err(err).
				Msg("failed to close kafka publisher")
		}
	}
}

