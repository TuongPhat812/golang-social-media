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

var _ contracts.OrderCancelledSubscriber = (*OrderCancelledSubscriber)(nil)

type OrderCancelledSubscriber struct {
	reader *kafka.Reader
}

func NewOrderCancelledSubscriber(brokers []string, groupID string) (*OrderCancelledSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.order_cancelled").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicOrderCancelled).
		Msg("order cancelled subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicOrderCancelled,
	})

	return &OrderCancelledSubscriber{reader: reader}, nil
}

func (s *OrderCancelledSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.order_cancelled").
		Info().
		Str("topic", events.TopicOrderCancelled).
		Msg("starting to consume OrderCancelled events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.order_cancelled").
				Info().
				Msg("stopping OrderCancelled subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.order_cancelled").
					Error().
					Err(err).
					Msg("failed to read message")
				continue
			}

			var event struct {
				OrderID    string
				UserID     string
				CancelledAt string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.order_cancelled").
					Error().
					Err(err).
					Msg("failed to unmarshal OrderCancelled event")
				continue
			}

			logger.Component("ecommerce.subscriber.order_cancelled").
				Info().
				Str("order_id", event.OrderID).
				Str("user_id", event.UserID).
				Msg("OrderCancelled event consumed")

			// TODO: Handle the event (e.g., send notification, refund, etc.)
		}
	}
}

func (s *OrderCancelledSubscriber) Close() error {
	return s.reader.Close()
}

