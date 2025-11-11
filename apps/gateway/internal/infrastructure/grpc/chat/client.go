package chat

import (
	"context"
	"log"
	"time"

	"github.com/myself/golang-social-media/pkg/config"
	chatcontract "github.com/myself/golang-social-media/pkg/contracts/chat"
	"github.com/myself/golang-social-media/pkg/grpcjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client chatcontract.ChatServiceClient
}

func NewClient(ctx context.Context) (*Client, error) {
	addr := config.GetEnv("CHAT_SERVICE_ADDR", "localhost:9000")

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcjson.Codec())),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("[gateway] connected to chat service at %s", addr)
	return &Client{conn: conn, client: chatcontract.NewChatServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CreateMessage(ctx context.Context, in *chatcontract.CreateMessageRequest) (*chatcontract.CreateMessageResponse, error) {
	return c.client.CreateMessage(ctx, in)
}
