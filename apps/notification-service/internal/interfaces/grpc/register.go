package grpc

import (
	"golang-social-media/apps/notification-service/internal/application/notifications"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func Register(server *grpc.Server, notificationService notifications.Service) {
	// TODO: register generated gRPC handlers when available.
	logger.Info().Msg("notification-service gRPC register invoked")
}
