package command

import (
	"context"
	"time"

	"golang-social-media/apps/gateway/internal/application/command/contracts"
	domain "golang-social-media/apps/gateway/internal/domain/message"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type createMessageCommand struct {
	client chatCommandClient
	log    *zerolog.Logger
}

type chatCommandClient interface {
	CreateMessage(ctx context.Context, in *chatv1.CreateMessageRequest, opts ...grpc.CallOption) (*chatv1.CreateMessageResponse, error)
}

func NewCreateMessageCommand(client chatCommandClient) contracts.CreateMessageCommand {
	return &createMessageCommand{
		client: client,
		log:    logger.Component("gateway.command.create_message"),
	}
}

func (c *createMessageCommand) Handle(ctx context.Context, senderID, receiverID, content string) (domain.Message, error) {
	resp, err := c.client.CreateMessage(ctx, &chatv1.CreateMessageRequest{
		SenderId:   senderID,
		ReceiverId: receiverID,
		Content:    content,
	})
	if err != nil {
		c.log.Error().
			Err(err).
			Str("sender_id", senderID).
			Str("receiver_id", receiverID).
			Msg("failed to call chat-service CreateMessage")
		return domain.Message{}, err
	}

	msg := resp.GetMessage()
	var createdAt time.Time
	if ts := msg.GetCreatedAt(); ts != nil {
		createdAt = ts.AsTime()
	}

	result := domain.Message{
		ID:         msg.GetId(),
		SenderID:   msg.GetSenderId(),
		ReceiverID: msg.GetReceiverId(),
		Content:    msg.GetContent(),
		CreatedAt:  createdAt,
	}

	c.log.Info().
		Str("message_id", result.ID).
		Str("sender_id", result.SenderID).
		Str("receiver_id", result.ReceiverID).
		Msg("message created")

	return result, nil
}
