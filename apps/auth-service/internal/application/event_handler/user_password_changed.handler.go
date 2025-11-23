package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/pkg/logger"
)

// UserPasswordChangedHandler handles UserPasswordChanged domain events
// Currently only logs the event. Can be extended to publish to Kafka if needed.
type UserPasswordChangedHandler struct {
	log *zerolog.Logger
}

// NewUserPasswordChangedHandler creates a new UserPasswordChangedHandler
func NewUserPasswordChangedHandler() *UserPasswordChangedHandler {
	return &UserPasswordChangedHandler{
		log: logger.Component("auth.event_handler.user_password_changed"),
	}
}

// Handle processes UserPasswordChanged domain events
func (h *UserPasswordChangedHandler) Handle(ctx context.Context, domainEvent user.DomainEvent) error {
	// Type assert to UserPasswordChangedEvent
	event, ok := domainEvent.(user.UserPasswordChangedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in UserPasswordChangedHandler")
		return nil // Ignore unexpected events
	}

	// Log the event (can be extended to publish to Kafka if other services need it)
	h.log.Info().
		Str("user_id", event.UserID).
		Str("updated_at", event.UpdatedAt).
		Msg("UserPasswordChanged event processed")

	// TODO: If other services need this event, publish to Kafka:
	// - Create UserPasswordChangedPayload in contracts
	// - Add PublishUserPasswordChanged to EventBrokerPublisher interface
	// - Implement in EventBrokerAdapter
	// - Publish to Kafka topic (e.g., "user.password.changed")
	// Note: Password change events are sensitive, consider security implications

	return nil
}

