package contracts

import (
	"context"
)

// EventBrokerPublisher publishes events to the event broker
// This is an abstraction over the infrastructure event bus (e.g., Kafka)
type EventBrokerPublisher interface {
	// PublishUserCreated publishes a user created event
	PublishUserCreated(ctx context.Context, payload UserCreatedPayload) error
}

// UserCreatedPayload represents the payload for user created event
type UserCreatedPayload struct {
	UserID    string
	Email     string
	Name      string
	CreatedAt string
}

