package publisher

import (
	"context"
	"time"

	"golang-social-media/apps/auth-service/internal/application/event_handler/contracts"
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

// PublishUserCreated publishes a user created event
func (a *EventBrokerAdapter) PublishUserCreated(ctx context.Context, payload contracts.UserCreatedPayload) error {
	// Transform application payload to infrastructure event
	createdAt, err := time.Parse(time.RFC3339, payload.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	kafkaEvent := events.UserCreated{
		ID:        payload.UserID,
		Email:     payload.Email,
		Name:      payload.Name,
		CreatedAt: createdAt,
	}

	return a.kafkaPublisher.PublishUserCreated(ctx, kafkaEvent)
}

