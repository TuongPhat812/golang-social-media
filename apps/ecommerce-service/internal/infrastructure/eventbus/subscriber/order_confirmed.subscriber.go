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

var _ contracts.OrderConfirmedSubscriber = (*OrderConfirmedSubscriber)(nil)

type OrderConfirmedSubscriber struct {
	reader *kafka.Reader
}

func NewOrderConfirmedSubscriber(brokers []string, groupID string) (*OrderConfirmedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.order_confirmed").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicOrderConfirmed).
		Msg("order confirmed subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicOrderConfirmed,
	})

	return &OrderConfirmedSubscriber{reader: reader}, nil
}

func (s *OrderConfirmedSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.order_confirmed").
		Info().
		Str("topic", events.TopicOrderConfirmed).
		Msg("starting to consume OrderConfirmed events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.order_confirmed").
				Info().
				Msg("stopping OrderConfirmed subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.order_confirmed").
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
				ConfirmedAt string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.order_confirmed").
					Error().
					Err(err).
					Msg("failed to unmarshal OrderConfirmed event")
				continue
			}

			logger.Component("ecommerce.subscriber.order_confirmed").
				Info().
				Str("order_id", event.OrderID).
				Str("user_id", event.UserID).
				Msg("OrderConfirmed event consumed")

			// TODO: Handle the event (e.g., send notification, trigger payment, etc.)
		}
	}
}

func (s *OrderConfirmedSubscriber) Close() error {
	return s.reader.Close()
}

