package messages

import (
	"context"
	"time"

	domain "golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/pkg/events"
)

type EventPublisher interface {
	PublishChatCreated(ctx context.Context, event events.ChatCreated) error
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error)
	SampleMessage() domain.Message
}

type service struct {
	publisher EventPublisher
}

func NewService(publisher EventPublisher) Service {
	return &service{publisher: publisher}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error) {
	createdAt := time.Now().UTC()
	msg := domain.Message{
		ID:         "msg-" + createdAt.Format("20060102150405"),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  createdAt,
	}

	event := events.ChatCreated{
		Message: events.ChatMessage{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt,
		},
		CreatedAt: createdAt,
	}
	if err := s.publisher.PublishChatCreated(ctx, event); err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *service) SampleMessage() domain.Message {
	now := time.Now().UTC()
	return domain.Message{ID: "sample", SenderID: "user-1", ReceiverID: "user-2", Content: "hello", CreatedAt: now}
}
