package contracts

import (
	"context"
)

// EventBrokerPublisher publishes events to the event broker
// This is an abstraction over the infrastructure event bus (e.g., Kafka)
type EventBrokerPublisher interface {
	// PublishNotificationCreated publishes a notification created event
	PublishNotificationCreated(ctx context.Context, payload NotificationCreatedPayload) error

	// PublishNotificationRead publishes a notification read event
	PublishNotificationRead(ctx context.Context, payload NotificationReadPayload) error
}

// NotificationCreatedPayload represents the payload for notification created event
type NotificationCreatedPayload struct {
	NotificationID string
	UserID         string
	Type           string
	Title          string
	Body           string
	Metadata       map[string]string
	CreatedAt      string
}

// NotificationReadPayload represents the payload for notification read event
type NotificationReadPayload struct {
	NotificationID string
	UserID         string
	ReadAt         string
}

