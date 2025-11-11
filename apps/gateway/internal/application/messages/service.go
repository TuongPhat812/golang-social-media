package messages

import (
	"context"

	domain "github.com/myself/golang-social-media/apps/gateway/internal/domain/message"
	"github.com/myself/golang-social-media/pkg/contracts/chat"
)

type ChatClient interface {
	CreateMessage(ctx context.Context, in *chat.CreateMessageRequest) (*chat.CreateMessageResponse, error)
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error)
}

type service struct {
	client ChatClient
}

func NewService(client ChatClient) Service {
	return &service{client: client}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error) {
	resp, err := s.client.CreateMessage(ctx, &chat.CreateMessageRequest{SenderID: senderID, ReceiverID: receiverID, Content: content})
	if err != nil {
		return domain.Message{}, err
	}
	return domain.Message{
		ID:         resp.Message.ID,
		SenderID:   resp.Message.SenderID,
		ReceiverID: resp.Message.ReceiverID,
		Content:    resp.Message.Content,
		CreatedAt:  resp.Message.CreatedAt,
	}, nil
}
