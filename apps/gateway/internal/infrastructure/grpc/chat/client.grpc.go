package chat

import (
	"context"
	"time"

	"golang-social-media/pkg/config"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client chatv1.ChatServiceClient
}

func NewClient(ctx context.Context) (*Client, error) {
	addr := config.GetEnv("CHAT_SERVICE_ADDR", "localhost:9000")

	// Increase timeout to 30 seconds to allow chat-service to be ready
	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	logger.Component("gateway.grpc").
		Info().
		Str("addr", addr).
		Msg("connecting to chat service")

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

	logger.Component("gateway.grpc").
		Info().
		Str("addr", addr).
		Msg("gateway connected to chat service")
	return &Client{conn: conn, client: chatv1.NewChatServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CreateMessage(ctx context.Context, in *chatv1.CreateMessageRequest, opts ...grpc.CallOption) (*chatv1.CreateMessageResponse, error) {
	return c.client.CreateMessage(ctx, in, opts...)
}
