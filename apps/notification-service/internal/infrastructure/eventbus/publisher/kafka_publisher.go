package publisher

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.NotificationPublisher = (*KafkaPublisher)(nil)

type KafkaPublisher struct {
	createdWriter *kafka.Writer
	readWriter    *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	createdWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicNotificationCreated,
		Balancer: &kafka.LeastBytes{},
	}

	readWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicNotificationRead,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Component("notification.publisher").
		Info().
		Strs("brokers", brokers).
		Msg("kafka publisher initialized")

	return &KafkaPublisher{
		createdWriter: createdWriter,
		readWriter:    readWriter,
	}, nil
}

func (p *KafkaPublisher) PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("notification.publisher").
			Error().
			Err(err).
			Msg("failed to marshal NotificationCreated event")
		return err
	}

	if err := p.createdWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Notification.ID),
		Value: payload,
	}); err != nil {
		logger.Component("notification.publisher").
			Error().
			Err(err).
			Msg("failed to publish NotificationCreated event")
		return err
	}

	logger.Component("notification.publisher").
		Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.Notification.ID).
		Msg("published NotificationCreated event")
	return nil
}

func (p *KafkaPublisher) PublishNotificationRead(ctx context.Context, event events.NotificationRead) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("notification.publisher").
			Error().
			Err(err).
			Msg("failed to marshal NotificationRead event")
		return err
	}

	if err := p.readWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.NotificationID),
		Value: payload,
	}); err != nil {
		logger.Component("notification.publisher").
			Error().
			Err(err).
			Msg("failed to publish NotificationRead event")
		return err
	}

	logger.Component("notification.publisher").
		Info().
		Str("topic", events.TopicNotificationRead).
		Str("notification_id", event.NotificationID).
		Str("user_id", event.UserID).
		Msg("published NotificationRead event")
	return nil
}

func (p *KafkaPublisher) Close() error {
	if err := p.createdWriter.Close(); err != nil {
		return err
	}
	return p.readWriter.Close()
}

