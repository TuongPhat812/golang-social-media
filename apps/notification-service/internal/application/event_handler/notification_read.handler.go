package event_handler

import (
	"context"

	"golang-social-media/apps/notification-service/internal/application/event_handler/contracts"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

type NotificationReadHandler struct {
	eventBroker contracts.EventBrokerPublisher
	log         *zerolog.Logger
}

func NewNotificationReadHandler(eventBroker contracts.EventBrokerPublisher) *NotificationReadHandler {
	return &NotificationReadHandler{
		eventBroker: eventBroker,
		log:         logger.Component("notification.event_handler.notification_read"),
	}
}

func (h *NotificationReadHandler) Handle(ctx context.Context, domainEvent notification.DomainEvent) error {
	// Type assert to NotificationReadEvent
	notificationReadEvent, ok := domainEvent.(notification.NotificationReadEvent)
	if !ok {
		h.log.Error().
			Str("event_type", domainEvent.Type()).
			Msg("unexpected event type in NotificationReadHandler")
		return nil // Ignore unexpected events
	}

	// Transform domain event to event broker payload
	payload := contracts.NotificationReadPayload{
		NotificationID: notificationReadEvent.NotificationID,
		UserID:         notificationReadEvent.UserID,
		ReadAt:         notificationReadEvent.ReadAt,
	}

	if err := h.eventBroker.PublishNotificationRead(ctx, payload); err != nil {
		h.log.Error().
			Err(err).
			Str("notification_id", notificationReadEvent.NotificationID).
			Msg("failed to publish NotificationRead event")
		return err
	}

	h.log.Info().
		Str("notification_id", notificationReadEvent.NotificationID).
		Str("user_id", notificationReadEvent.UserID).
		Msg("NotificationRead event published")

	return nil
}
