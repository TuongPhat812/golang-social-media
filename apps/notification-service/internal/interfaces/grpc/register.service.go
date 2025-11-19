package grpc

import (
	bootstrap "golang-social-media/apps/notification-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

// RegisterServices registers all gRPC services with the server
func RegisterServices(server *grpc.Server, deps *bootstrap.Dependencies) {
	// TODO: Register gRPC handlers when available
	// Example:
	// notificationpb.RegisterNotificationServiceServer(server, NewNotificationServiceHandler(deps))
	
	logger.Component("notification.grpc").
		Info().
		Msg("gRPC services registered")
}
