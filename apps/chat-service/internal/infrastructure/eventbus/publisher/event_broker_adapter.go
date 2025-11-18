package publisher

import (
	"context"
	"time"

	"golang-social-media/apps/chat-service/internal/application/event_handler/contracts"
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

// PublishMessageCreated publishes a message created event
func (a *EventBrokerAdapter) PublishMessageCreated(ctx context.Context, payload contracts.MessageCreatedPayload) error {
	// Transform application payload to infrastructure event
	createdAt, err := time.Parse(time.RFC3339, payload.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	kafkaEvent := events.ChatCreated{
		Message: events.ChatMessage{
			ID:         payload.MessageID,
			SenderID:   payload.SenderID,
			ReceiverID: payload.ReceiverID,
			Content:    payload.Content,
			CreatedAt:  createdAt,
		},
		CreatedAt: createdAt,
	}

	return a.kafkaPublisher.PublishChatCreated(ctx, kafkaEvent)
}
