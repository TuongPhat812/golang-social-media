package eventbus

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicNotificationCreated,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Info().
		Strs("brokers", brokers).
		Msg("notification-service kafka publisher initialized")

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Error().Err(err).Msg("notification-service failed to marshal NotificationCreated event")
		return err
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Notification.ID),
		Value: payload,
	}); err != nil {
		logger.Error().Err(err).Msg("notification-service failed to publish NotificationCreated event")
		return err
	}

	logger.Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.Notification.ID).
		Msg("notification-service published NotificationCreated event")
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
