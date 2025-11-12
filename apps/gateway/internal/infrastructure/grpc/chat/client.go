package chat

import (
	"context"
	"log"
	"time"

	"golang-social-media/pkg/config"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client chatv1.ChatServiceClient
}

func NewClient(ctx context.Context) (*Client, error) {
	addr := config.GetEnv("CHAT_SERVICE_ADDR", "localhost:9000")

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("[gateway] connected to chat service at %s", addr)
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
