package grpc

import (
	"fmt"
	"net"

	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Start starts the gRPC server
func Start(addr string, registerFunc func(*grpc.Server)) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Create gRPC server with optimizations
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * 60, // 15 minutes
			MaxConnectionAge:      30 * 60, // 30 minutes
			MaxConnectionAgeGrace: 5 * 60,  // 5 minutes
			Time:                  5 * 60,  // 5 minutes
			Timeout:               20,      // 20 seconds
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * 60, // 5 minutes
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
	)

	// Register services
	registerFunc(server)

	logger.Component("auth.grpc").
		Info().
		Str("addr", addr).
		Msg("auth gRPC server starting")

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

