package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/pkg/logger"
)

type UserCreatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewUserCreatedHandler(eventBroker contracts.EventBrokerPublisher) *UserCreatedHandler {
	return &UserCreatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("auth.event_handler.user_created"),
	}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, domainEvent user.DomainEvent) error {
	// Type assert to UserCreatedEvent
	userCreatedEvent, ok := domainEvent.(user.UserCreatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in UserCreatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.UserCreatedPayload{
		UserID:    userCreatedEvent.UserID,
		Email:     userCreatedEvent.Email,
		Name:      userCreatedEvent.Name,
		CreatedAt: userCreatedEvent.CreatedAt,
	}

	if err := h.eventBroker.PublishUserCreated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("user_id", userCreatedEvent.UserID).
			Msg("failed to publish UserCreated event")
		return err
	}

	h.log.Info().
		Str("user_id", userCreatedEvent.UserID).
		Str("email", userCreatedEvent.Email).
		Msg("UserCreated event published")

	return nil
}

