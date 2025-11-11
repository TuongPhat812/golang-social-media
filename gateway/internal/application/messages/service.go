package messages

import (
	"context"

	"github.com/myself/golang-social-media/common/contracts/chat"
	"github.com/myself/golang-social-media/common/domain/message"
)

type ChatClient interface {
	CreateMessage(ctx context.Context, in *chat.CreateMessageRequest) (*chat.CreateMessageResponse, error)
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (message.Message, error)
}

type service struct {
	client ChatClient
}

func NewService(client ChatClient) Service {
	return &service{client: client}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (message.Message, error) {
	resp, err := s.client.CreateMessage(ctx, &chat.CreateMessageRequest{SenderID: senderID, ReceiverID: receiverID, Content: content})
	if err != nil {
		return message.Message{}, err
	}
	return resp.Message, nil
}
