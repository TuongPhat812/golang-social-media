package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/segmentio/kafka-go"
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

		// Batching Configuration - optimized for low latency (user creation is critical)
		BatchSize:    50,                    // Smaller batch for lower latency
		BatchBytes:   1048576,               // 1MB max batch size
		BatchTimeout: 10 * time.Millisecond, // Flush every 10ms

		// Reliability
		RequiredAcks: kafka.RequireOne, // Wait for leader ack
		MaxAttempts:  10,               // Retry up to 10 times

		// Timeouts
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		// Backoff for retries
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,

		// Performance
		Async: true, // Non-blocking writes

		// Compression - Snappy for JSON payloads (typically 200-500 bytes)
		Compression: kafka.Snappy,
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
