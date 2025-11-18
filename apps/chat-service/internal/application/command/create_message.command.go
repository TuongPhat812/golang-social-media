package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang-social-media/apps/chat-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/chat-service/internal/application/event_dispatcher"
	"golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	"golang-social-media/pkg/logger"
)

var _ contracts.CreateMessageCommand = (*createMessageCommand)(nil)

type createMessageCommand struct {
	repo           *persistence.MessageRepository
	eventDispatcher *event_dispatcher.Dispatcher
	log            *zerolog.Logger
}

func NewCreateMessageCommand(
	repo *persistence.MessageRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CreateMessageCommand {
	return &createMessageCommand{
		repo:           repo,
		eventDispatcher: eventDispatcher,
		log:            logger.Component("chat.command.create_message"),
	}
}

func (c *createMessageCommand) Execute(ctx context.Context, req contracts.CreateMessageCommandRequest) (message.Message, error) {
	createdAt := time.Now().UTC()
	messageModel := message.Message{
		ID:         uuid.NewString(),
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		CreatedAt:  createdAt,
	}

	// Validate business rules before persisting or publishing
	if err := messageModel.Validate(); err != nil {
		c.log.Error().
			Err(err).
			Str("sender_id", req.SenderID).
			Str("receiver_id", req.ReceiverID).
			Msg("message validation failed")
		return message.Message{}, err
	}

	// Domain logic: create message (this adds domain events internally)
	messageModel.Create()

	// Save domain events BEFORE persisting (repository might overwrite the message)
	domainEvents := messageModel.Events()

	// Persist to database
	if err := c.repo.Create(ctx, &messageModel); err != nil {
		c.log.Error().
			Err(err).
			Str("sender_id", req.SenderID).
			Str("receiver_id", req.ReceiverID).
			Msg("failed to persist message")
		return message.Message{}, err
	}

	// Dispatch domain events AFTER successful persistence
	messageModel.ClearEvents() // Clear events after dispatch

	c.log.Info().
		Int("event_count", len(domainEvents)).
		Str("message_id", messageModel.ID).
		Msg("dispatching domain events")

	for _, domainEvent := range domainEvents {
		c.log.Info().
			Str("event_type", domainEvent.Type()).
			Str("message_id", messageModel.ID).
			Msg("dispatching domain event")

		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			// Log error but don't fail the command
			// Events can be retried via outbox pattern in production
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("message_id", messageModel.ID).
				Msg("failed to dispatch domain event")
		} else {
			c.log.Debug().
				Str("event_type", domainEvent.Type()).
				Str("message_id", messageModel.ID).
				Msg("domain event dispatched successfully")
		}
	}

	c.log.Info().
		Str("message_id", messageModel.ID).
		Str("sender_id", messageModel.SenderID).
		Str("receiver_id", messageModel.ReceiverID).
		Msg("message created")

	return messageModel, nil
}

