package messages

import (
	"context"
	"time"

	domain "golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishChatCreated(ctx context.Context, event events.ChatCreated) error
}

type Service interface {
	CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error)
}

type service struct {
	repository     Repository
	eventPublisher EventPublisher
}

func NewService(repository Repository, eventPublisher EventPublisher) Service {
	return &service{
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}

func (s *service) CreateMessage(ctx context.Context, senderID, receiverID, content string) (domain.Message, error) {
	createdAt := time.Now().UTC()
	msg := domain.Message{
		ID:         uuid.NewString(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  createdAt,
	}

	if err := s.repository.Create(ctx, &msg); err != nil {
		logger.Error().Err(err).Msg("failed to persist chat message")
		return msg, err
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
	if err := s.eventPublisher.PublishChatCreated(ctx, event); err != nil {
		logger.Error().Err(err).Msg("failed to publish ChatCreated event")
		return msg, err
	}

	logger.Info().
		Str("message_id", msg.ID).
		Str("sender_id", msg.SenderID).
		Str("receiver_id", msg.ReceiverID).
		Msg("chat message created and event published")

	return msg, nil
}
