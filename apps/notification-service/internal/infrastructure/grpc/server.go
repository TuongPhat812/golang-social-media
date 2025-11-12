package grpc

import (
	"net"

	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func Start(addr string, register func(*grpc.Server)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	if register != nil {
		register(server)
	}

	logger.Info().Str("addr", addr).Msg("notification-service gRPC server starting")
	return server.Serve(listener)
}
