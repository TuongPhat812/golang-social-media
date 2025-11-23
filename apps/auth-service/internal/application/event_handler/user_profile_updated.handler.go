package event_handler

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/pkg/logger"
)

// UserProfileUpdatedHandler handles UserProfileUpdated domain events
// Currently only logs the event. Can be extended to publish to Kafka if needed.
type UserProfileUpdatedHandler struct {
	log *zerolog.Logger
}

// NewUserProfileUpdatedHandler creates a new UserProfileUpdatedHandler
func NewUserProfileUpdatedHandler() *UserProfileUpdatedHandler {
	return &UserProfileUpdatedHandler{
		log: logger.Component("auth.event_handler.user_profile_updated"),
	}
}

// Handle processes UserProfileUpdated domain events
func (h *UserProfileUpdatedHandler) Handle(ctx context.Context, domainEvent user.DomainEvent) error {
	// Type assert to UserProfileUpdatedEvent
	event, ok := domainEvent.(user.UserProfileUpdatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in UserProfileUpdatedHandler")
		return nil // Ignore unexpected events
	}

	// Log the event (can be extended to publish to Kafka if other services need it)
	h.log.Info().
		Str("user_id", event.UserID).
		Str("old_name", event.OldName).
		Str("new_name", event.NewName).
		Str("updated_at", event.UpdatedAt).
		Msg("UserProfileUpdated event processed")

	// TODO: If other services need this event, publish to Kafka:
	// - Create UserProfileUpdatedPayload in contracts
	// - Add PublishUserProfileUpdated to EventBrokerPublisher interface
	// - Implement in EventBrokerAdapter
	// - Publish to Kafka topic (e.g., "user.profile.updated")

	return nil
}

