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

	logger.Component("notification.grpc").
		Info().
		Str("addr", addr).
		Msg("gRPC server starting")
	return server.Serve(listener)
}
