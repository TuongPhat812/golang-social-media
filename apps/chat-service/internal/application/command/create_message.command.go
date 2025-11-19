package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/chat-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/chat-service/internal/application/event_dispatcher"
	"golang-social-media/apps/chat-service/internal/domain/factories"
	"golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	"golang-social-media/pkg/logger"
)

var _ contracts.CreateMessageCommand = (*createMessageCommand)(nil)

type createMessageCommand struct {
	repo            *persistence.MessageRepository
	messageFactory  *factories.MessageFactory
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewCreateMessageCommand(
	repo *persistence.MessageRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CreateMessageCommand {
	return &createMessageCommand{
		repo:            repo,
		messageFactory:  factories.NewMessageFactory(),
		eventDispatcher: eventDispatcher,
		log:             logger.Component("chat.command.create_message"),
	}
}

func (c *createMessageCommand) Execute(ctx context.Context, req contracts.CreateMessageCommandRequest) (message.Message, error) {
	startTime := time.Now()

	// Use factory to create message
	modelStart := time.Now()
	messageModel, err := c.messageFactory.CreateMessage(req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		modelDuration := time.Since(modelStart)
		totalDuration := time.Since(startTime)
		c.log.Error().
			Err(err).
			Str("sender_id", req.SenderID).
			Str("receiver_id", req.ReceiverID).
			Dur("model_create_ms", modelDuration).
			Dur("total_ms", totalDuration).
			Msg("failed to create message using factory")
		return message.Message{}, err
	}
	modelDuration := time.Since(modelStart)

	// Save domain events BEFORE persisting (repository might overwrite the message)
	domainEvents := messageModel.Events()

	// Persist to database
	dbStart := time.Now()
	if err := c.repo.Create(ctx, messageModel); err != nil {
		dbDuration := time.Since(dbStart)
		totalDuration := time.Since(startTime)
		c.log.Error().
			Err(err).
			Str("sender_id", req.SenderID).
			Str("receiver_id", req.ReceiverID).
			Dur("model_create_ms", modelDuration).
			Dur("db_persist_ms", dbDuration).
			Dur("total_ms", totalDuration).
			Msg("failed to persist message")
		return message.Message{}, err
	}
	dbDuration := time.Since(dbStart)

	// Dispatch domain events AFTER successful persistence
	messageModel.ClearEvents() // Clear events after dispatch

	dispatchStart := time.Now()
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
	dispatchDuration := time.Since(dispatchStart)

	totalDuration := time.Since(startTime)

	c.log.Info().
		Str("message_id", messageModel.ID).
		Str("sender_id", messageModel.SenderID).
		Str("receiver_id", messageModel.ReceiverID).
		Dur("model_create_ms", modelDuration).
		Dur("db_persist_ms", dbDuration).
		Dur("event_dispatch_ms", dispatchDuration).
		Dur("total_ms", totalDuration).
		Msg("message created")

	return *messageModel, nil
}

