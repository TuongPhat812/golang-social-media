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

var _ contracts.ProductStockUpdatedSubscriber = (*ProductStockUpdatedSubscriber)(nil)

type ProductStockUpdatedSubscriber struct {
	reader *kafka.Reader
}

func NewProductStockUpdatedSubscriber(brokers []string, groupID string) (*ProductStockUpdatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.product_stock_updated").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicProductStockUpdated).
		Msg("product stock updated subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicProductStockUpdated,
	})

	return &ProductStockUpdatedSubscriber{reader: reader}, nil
}

func (s *ProductStockUpdatedSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.product_stock_updated").
		Info().
		Str("topic", events.TopicProductStockUpdated).
		Msg("starting to consume ProductStockUpdated events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.product_stock_updated").
				Info().
				Msg("stopping ProductStockUpdated subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.product_stock_updated").
					Error().
					Err(err).
					Msg("failed to read message")
				continue
			}

			var event struct {
				ProductID string
				OldStock  int
				NewStock  int
				UpdatedAt string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.product_stock_updated").
					Error().
					Err(err).
					Msg("failed to unmarshal ProductStockUpdated event")
				continue
			}

			logger.Component("ecommerce.subscriber.product_stock_updated").
				Info().
				Str("product_id", event.ProductID).
				Int("old_stock", event.OldStock).
				Int("new_stock", event.NewStock).
				Msg("ProductStockUpdated event consumed")

			// TODO: Handle the event (e.g., update cache, send notification, etc.)
		}
	}
}

func (s *ProductStockUpdatedSubscriber) Close() error {
	return s.reader.Close()
}

