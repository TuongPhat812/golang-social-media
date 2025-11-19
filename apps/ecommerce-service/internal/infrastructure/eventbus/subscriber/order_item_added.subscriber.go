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

var _ contracts.OrderItemAddedSubscriber = (*OrderItemAddedSubscriber)(nil)

type OrderItemAddedSubscriber struct {
	reader *kafka.Reader
}

func NewOrderItemAddedSubscriber(brokers []string, groupID string) (*OrderItemAddedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.order_item_added").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicOrderItemAdded).
		Msg("order item added subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicOrderItemAdded,
	})

	return &OrderItemAddedSubscriber{reader: reader}, nil
}

func (s *OrderItemAddedSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.order_item_added").
		Info().
		Str("topic", events.TopicOrderItemAdded).
		Msg("starting to consume OrderItemAdded events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.order_item_added").
				Info().
				Msg("stopping OrderItemAdded subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.order_item_added").
					Error().
					Err(err).
					Msg("failed to read message")
				continue
			}

			var event struct {
				OrderID   string
				ProductID string
				Quantity  int
				UnitPrice float64
				SubTotal  float64
				UpdatedAt string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.order_item_added").
					Error().
					Err(err).
					Msg("failed to unmarshal OrderItemAdded event")
				continue
			}

			logger.Component("ecommerce.subscriber.order_item_added").
				Info().
				Str("order_id", event.OrderID).
				Str("product_id", event.ProductID).
				Int("quantity", event.Quantity).
				Msg("OrderItemAdded event consumed")

			// TODO: Handle the event (e.g., update cache, send notification, etc.)
		}
	}
}

func (s *OrderItemAddedSubscriber) Close() error {
	return s.reader.Close()
}

