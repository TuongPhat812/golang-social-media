package grpc

import (
	"log"

	"golang-social-media/apps/notification-service/internal/application/notifications"
	"google.golang.org/grpc"
)

func Register(server *grpc.Server, notificationService notifications.Service) {
	// TODO: register generated gRPC handlers when available.
	log.Printf("grpc.Register invoked with notification service: %+v", notificationService.SampleUser())
}
