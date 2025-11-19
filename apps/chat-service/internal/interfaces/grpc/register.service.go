package grpc

import (
	bootstrap "golang-social-media/apps/chat-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/logger"

	"google.golang.org/grpc"
)

// RegisterServices registers all gRPC services with the server
func RegisterServices(server *grpc.Server, deps *bootstrap.Dependencies) {
	// TODO: Register gRPC handlers when available
	// Example:
	// chatpb.RegisterChatServiceServer(server, NewChatServiceHandler(deps))

	logger.Component("chat.grpc").
		Info().
		Msg("gRPC services registered")
}
