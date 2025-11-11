package grpc

import (
	"log"
	"net"

	"github.com/myself/golang-social-media/common/grpcjson"
	"google.golang.org/grpc"
)

func Start(addr string, register func(*grpc.Server)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.ForceServerCodec(grpcjson.Codec()))
	if register != nil {
		register(server)
	}

	log.Printf("chat service gRPC server starting on %s", addr)
	return server.Serve(listener)
}
