package messages

import (
	"context"
	"time"

	domain "golang-social-media/apps/gateway/internal/domain/message"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type ChatClient interface {
	CreateMessage(ctx context.Context, in *chatv1.CreateMessageRequest, opts ...grpc.CallOption) (*chatv1.CreateMessageResponse, error)
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error)
}

type service struct {
	client ChatClient
	log    *zerolog.Logger
}

func NewService(client ChatClient, log *zerolog.Logger) Service {
	if log == nil {
		log = logger.Component("gateway.messages")
	}
	return &service{client: client, log: log}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error) {
	resp, err := s.client.CreateMessage(ctx, &chatv1.CreateMessageRequest{SenderId: senderID, ReceiverId: receiverID, Content: content})
	if err != nil {
		s.log.Error().
			Err(err).
			Str("sender_id", senderID).
			Str("receiver_id", receiverID).
			Msg("failed to invoke chat service CreateMessage")
		return domain.Message{}, err
	}
	msg := resp.GetMessage()
	var createdAt time.Time
	if ts := msg.GetCreatedAt(); ts != nil {
		createdAt = ts.AsTime()
	}
	return domain.Message{
		ID:         msg.GetId(),
		SenderID:   msg.GetSenderId(),
		ReceiverID: msg.GetReceiverId(),
		Content:    msg.GetContent(),
		CreatedAt:  createdAt,
	}, nil
}
