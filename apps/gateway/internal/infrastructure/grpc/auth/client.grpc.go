package auth

import (
	"context"
	"time"

	"golang-social-media/pkg/config"
	authv1 "golang-social-media/pkg/gen/auth/v1"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewClient(ctx context.Context) (*Client, error) {
	addr := config.GetEnv("AUTH_SERVICE_ADDR", "localhost:9100")

	// Increase timeout to 30 seconds to allow auth-service to be ready
	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logger.Component("gateway.grpc.auth").
		Info().
		Str("addr", addr).
		Msg("connecting to auth service")

	conn, err := grpc.DialContext(
		dialCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithInitialWindowSize(65535),       // Increase initial window size
		grpc.WithInitialConnWindowSize(1048576), // 1MB initial connection window
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB max receive message size
			grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB max send message size
		),
	)
	if err != nil {
		return nil, err
	}

	logger.Component("gateway.grpc.auth").
		Info().
		Str("addr", addr).
		Msg("gateway connected to auth service")
	return &Client{conn: conn, client: authv1.NewAuthServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	req := &authv1.ValidateTokenRequest{
		Token: token,
	}
	return c.client.ValidateToken(ctx, req)
}

