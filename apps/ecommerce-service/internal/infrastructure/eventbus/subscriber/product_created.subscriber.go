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

var _ contracts.ProductCreatedSubscriber = (*ProductCreatedSubscriber)(nil)

type ProductCreatedSubscriber struct {
	reader *kafka.Reader
}

func NewProductCreatedSubscriber(brokers []string, groupID string) (*ProductCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("ecommerce.subscriber.product_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicProductCreated).
		Msg("product created subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicProductCreated,
	})

	return &ProductCreatedSubscriber{reader: reader}, nil
}

func (s *ProductCreatedSubscriber) Consume(ctx context.Context) {
	logger.Component("ecommerce.subscriber.product_created").
		Info().
		Str("topic", events.TopicProductCreated).
		Msg("starting to consume ProductCreated events")

	for {
		select {
		case <-ctx.Done():
			logger.Component("ecommerce.subscriber.product_created").
				Info().
				Msg("stopping ProductCreated subscriber")
			return
		default:
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				logger.Component("ecommerce.subscriber.product_created").
					Error().
					Err(err).
					Msg("failed to read message")
				continue
			}

			var event struct {
				ProductID   string
				Name        string
				Description string
				Price       float64
				Stock       int
				CreatedAt   string
			}

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("ecommerce.subscriber.product_created").
					Error().
					Err(err).
					Msg("failed to unmarshal ProductCreated event")
				continue
			}

			logger.Component("ecommerce.subscriber.product_created").
				Info().
				Str("product_id", event.ProductID).
				Str("name", event.Name).
				Msg("ProductCreated event consumed")

			// TODO: Handle the event (e.g., update cache, send notification, etc.)
		}
	}
}

func (s *ProductCreatedSubscriber) Close() error {
	return s.reader.Close()
}

