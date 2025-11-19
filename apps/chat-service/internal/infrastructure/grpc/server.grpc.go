package grpc

import (
	"net"
	"os"

	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func Start(addr string, register func(*grpc.Server)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Initialize error transformer
	devMode := os.Getenv("ENV") == "development"
	transformer := errors.NewTransformer(devMode)

	// Create server with error interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(errors.GRPCErrorInterceptor(transformer)),
	)
	if register != nil {
		register(server)
	}

	logger.Component("chat.grpc").
		Info().
		Str("addr", addr).
		Msg("chat-service gRPC server starting")
	return server.Serve(listener)
}
