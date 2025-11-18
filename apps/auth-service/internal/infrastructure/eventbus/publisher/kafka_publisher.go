package publisher

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.UserPublisher = (*KafkaPublisher)(nil)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicUserCreated,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Component("auth.publisher").
		Info().
		Strs("brokers", brokers).
		Msg("kafka publisher initialized")

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishUserCreated(ctx context.Context, event events.UserCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("auth.publisher").
			Error().
			Err(err).
			Msg("failed to marshal UserCreated event")
		return err
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: payload,
	}); err != nil {
		logger.Component("auth.publisher").
			Error().
			Err(err).
			Msg("failed to publish UserCreated event")
		return err
	}

	logger.Component("auth.publisher").
		Info().
		Str("topic", events.TopicUserCreated).
		Str("user_id", event.ID).
		Msg("published UserCreated event")
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

