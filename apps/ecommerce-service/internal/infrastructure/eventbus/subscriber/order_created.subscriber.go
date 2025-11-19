package subscriber

import (
	"context"
	"encoding/json"
	"errors"

	"golang-social-media/apps/ecommerce-service/internal/infrastructure/eventbus/subscriber/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/segmentio/kafka-go"
)

var _ contracts.OrderCreatedSubscriber = (*OrderCreatedSubscriber)(nil)

type OrderCreatedSubscriber struct {
	reader *kafka.Reader
}

func NewOrderCreatedSubscriber(brokers []string, groupID string) (*OrderCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.order_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicOrderCreated).
		Msg("order created subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicOrderCreated,
	})

	return &OrderCreatedSubscriber{reader: reader}, nil
}

func (s *OrderCreatedSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.order_created").
		Info().
		Str("topic", events.TopicOrderCreated).
		Msg("starting to consume OrderCreated events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.order_created").
				Info().
				Msg("stopping OrderCreated subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.order_created").
					Error().
					Err(err).
					Msg("failed to read message")
				continue
			}

			var event struct {
				OrderID     string
				UserID      string
				TotalAmount float64
				ItemCount   int
				CreatedAt   string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.order_created").
					Error().
					Err(err).
					Msg("failed to unmarshal OrderCreated event")
				continue
			}

			logger.Component("ecommerce.subscriber.order_created").
				Info().
				Str("order_id", event.OrderID).
				Str("user_id", event.UserID).
				Msg("OrderCreated event consumed")

			// TODO: Handle the event (e.g., update cache, send notification, etc.)
		}
	}
}

func (s *OrderCreatedSubscriber) Close() error {
	return s.reader.Close()
}

