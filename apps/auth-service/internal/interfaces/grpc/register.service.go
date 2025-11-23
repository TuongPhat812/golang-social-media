package grpc

import (
	bootstrap "golang-social-media/apps/auth-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/logger"

	"google.golang.org/grpc"
)

// RegisterServices registers all gRPC services with the server
func RegisterServices(server *grpc.Server, deps *bootstrap.Dependencies) {
	logger.Component("auth.grpc").
		Info().
		Msg("gRPC services registered")
}

