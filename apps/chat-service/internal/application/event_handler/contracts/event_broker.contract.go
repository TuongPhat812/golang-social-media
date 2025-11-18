package contracts

import (
	"context"
)

// EventBrokerPublisher publishes events to the event broker
// This is an abstraction over the infrastructure event bus (e.g., Kafka)
type EventBrokerPublisher interface {
	// PublishMessageCreated publishes a message created event
	PublishMessageCreated(ctx context.Context, payload MessageCreatedPayload) error
}

// MessageCreatedPayload represents the payload for message created event
type MessageCreatedPayload struct {
	MessageID  string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  string
}

