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

	// Create server with error interceptor and performance optimizations
	server := grpc.NewServer(
		grpc.UnaryInterceptor(errors.GRPCErrorInterceptor(transformer)),
		grpc.MaxConcurrentStreams(10000),        // Allow high concurrency
		grpc.InitialWindowSize(65535),          // Increase initial window size
		grpc.InitialConnWindowSize(1048576),    // 1MB initial connection window
		grpc.MaxRecvMsgSize(4*1024*1024),       // 4MB max receive message size
		grpc.MaxSendMsgSize(4*1024*1024),       // 4MB max send message size
		grpc.NumStreamWorkers(10),              // Use multiple workers for streams
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
