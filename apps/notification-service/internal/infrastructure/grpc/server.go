package grpc

import (
	"log"
	"net"

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

	log.Printf("notification service gRPC server starting on %s", addr)
	return server.Serve(listener)
}
