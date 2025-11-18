package publisher

import (
	"context"
	"time"

	"golang-social-media/apps/notification-service/internal/application/event_handler/contracts"
	"golang-social-media/pkg/events"
)

// EventBrokerAdapter adapts the infrastructure Kafka publisher to the application event broker interface
type EventBrokerAdapter struct {
	kafkaPublisher *KafkaPublisher
}

// NewEventBrokerAdapter creates a new event broker adapter
func NewEventBrokerAdapter(kafkaPublisher *KafkaPublisher) *EventBrokerAdapter {
	return &EventBrokerAdapter{
		kafkaPublisher: kafkaPublisher,
	}
}

// PublishNotificationCreated publishes a notification created event
func (a *EventBrokerAdapter) PublishNotificationCreated(ctx context.Context, payload contracts.NotificationCreatedPayload) error {
	// Transform application payload to infrastructure event
	createdAt, err := time.Parse(time.RFC3339, payload.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	kafkaEvent := events.NotificationCreated{
		Notification: events.Notification{
			ID:        payload.NotificationID,
			UserID:    payload.UserID,
			Type:      payload.Type,
			Title:     payload.Title,
			Body:      payload.Body,
			Metadata:  payload.Metadata,
			CreatedAt: createdAt,
		},
	}

	return a.kafkaPublisher.PublishNotificationCreated(ctx, kafkaEvent)
}

// PublishNotificationRead publishes a notification read event
func (a *EventBrokerAdapter) PublishNotificationRead(ctx context.Context, payload contracts.NotificationReadPayload) error {
	// Transform application payload to infrastructure event
	readAt, err := time.Parse(time.RFC3339, payload.ReadAt)
	if err != nil {
		readAt = time.Now()
	}

	kafkaEvent := events.NotificationRead{
		NotificationID: payload.NotificationID,
		UserID:         payload.UserID,
		ReadAt:         readAt,
	}

	return a.kafkaPublisher.PublishNotificationRead(ctx, kafkaEvent)
}

