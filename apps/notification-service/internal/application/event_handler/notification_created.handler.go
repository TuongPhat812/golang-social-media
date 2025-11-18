package event_handler

import (
	"context"

	"golang-social-media/apps/notification-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

type NotificationCreatedHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewNotificationCreatedHandler(eventBroker contracts.EventBrokerPublisher) *NotificationCreatedHandler {
	return &NotificationCreatedHandler{
		eventBroker: eventBroker,
		log:         logger.Component("notification.event_handler.notification_created"),
	}
}

func (h *NotificationCreatedHandler) Handle(ctx context.Context, domainEvent notification.DomainEvent) error {
	// Type assert to NotificationCreatedEvent
	notificationCreatedEvent, ok := domainEvent.(notification.NotificationCreatedEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in NotificationCreatedHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.NotificationCreatedPayload{
		NotificationID: notificationCreatedEvent.NotificationID,
		UserID:         notificationCreatedEvent.UserID,
		Type:           string(notificationCreatedEvent.NotificationType),
		Title:          notificationCreatedEvent.Title,
		Body:           notificationCreatedEvent.Body,
		Metadata:       notificationCreatedEvent.Metadata,
		CreatedAt:      notificationCreatedEvent.CreatedAt,
	}

	if err := h.eventBroker.PublishNotificationCreated(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("notification_id", notificationCreatedEvent.NotificationID).
			Msg("failed to publish NotificationCreated event")
		return err
	}

	h.log.Info().
		Str("notification_id", notificationCreatedEvent.NotificationID).
		Str("user_id", notificationCreatedEvent.UserID).
		Msg("NotificationCreated event published")

	return nil
}
