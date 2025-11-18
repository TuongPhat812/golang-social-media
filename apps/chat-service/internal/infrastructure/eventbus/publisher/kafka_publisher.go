package publisher

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/chat-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.ChatPublisher = (*KafkaPublisher)(nil)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicChatCreated,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Component("chat.publisher").
		Info().
		Strs("brokers", brokers).
		Msg("kafka publisher initialized")

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishChatCreated(ctx context.Context, event events.ChatCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("chat.publisher").
			Error().
			Err(err).
			Msg("failed to marshal ChatCreated event")
		return err
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Message.ID),
		Value: payload,
	}); err != nil {
		logger.Component("chat.publisher").
			Error().
			Err(err).
			Msg("failed to publish ChatCreated event")
		return err
	}

	logger.Component("chat.publisher").
		Info().
		Str("topic", events.TopicChatCreated).
		Str("message_id", event.Message.ID).
		Msg("published ChatCreated event")
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

