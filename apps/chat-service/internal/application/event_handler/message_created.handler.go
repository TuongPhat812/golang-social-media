package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/chat-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/pkg/logger"
)

type MessageCreatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewMessageCreatedHandler(eventBroker contracts.EventBrokerPublisher) *MessageCreatedHandler {
	return &MessageCreatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("chat.event_handler.message_created"),
	}
}

func (h *MessageCreatedHandler) Handle(ctx context.Context, domainEvent message.DomainEvent) error {
	// Type assert to MessageCreatedEvent
	messageCreatedEvent, ok := domainEvent.(message.MessageCreatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in MessageCreatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.MessageCreatedPayload{
		MessageID:  messageCreatedEvent.MessageID,
		SenderID:   messageCreatedEvent.SenderID,
		ReceiverID: messageCreatedEvent.ReceiverID,
		Content:    messageCreatedEvent.Content,
		CreatedAt:  messageCreatedEvent.CreatedAt,
	}

	if err := h.eventBroker.PublishMessageCreated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("message_id", messageCreatedEvent.MessageID).
			Msg("failed to publish MessageCreated event")
		return err
	}

	h.log.Info().
		Str("message_id", messageCreatedEvent.MessageID).
		Str("sender_id", messageCreatedEvent.SenderID).
		Str("receiver_id", messageCreatedEvent.ReceiverID).
		Msg("MessageCreated event published")

	return nil
}

