package messages

import (
	"context"
	"time"

	"github.com/myself/golang-social-media/common/domain/message"
	"github.com/myself/golang-social-media/common/events"
)

type EventPublisher interface {
	PublishChatCreated(ctx context.Context, event events.ChatCreated) error
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (message.Message, error)
	SampleMessage() message.Message
}

type service struct {
	publisher EventPublisher
}

func NewService(publisher EventPublisher) Service {
	return &service{publisher: publisher}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (message.Message, error) {
	createdAt := time.Now().UTC()
	msg := message.Message{
		ID:         "msg-" + createdAt.Format("20060102150405"),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  createdAt,
	}

	event := events.ChatCreated{Message: msg, CreatedAt: createdAt}
	if err := s.publisher.PublishChatCreated(ctx, event); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SampleMessage() message.Message {
	now := time.Now().UTC()
	return message.Message{ID: "sample", SenderID: "user-1", ReceiverID: "user-2", Content: "hello", CreatedAt: now}
}
