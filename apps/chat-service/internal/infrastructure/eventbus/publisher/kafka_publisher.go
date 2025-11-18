package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang-social-media/apps/chat-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/segmentio/kafka-go"
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

		// Batching Configuration - optimized for throughput with low latency
		BatchSize:    100,                   // Batch up to 100 messages for better throughput
		BatchBytes:   1048576,               // 1MB max batch size
		BatchTimeout: 10 * time.Millisecond, // Flush every 10ms if batch incomplete

		// Reliability
		RequiredAcks: kafka.RequireOne, // Wait for leader ack (good balance)
		MaxAttempts:  10,               // Retry up to 10 times

		// Timeouts
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		// Backoff for retries
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,

		// Performance
		Async: true, // Non-blocking writes

		// Compression - Snappy provides best balance of speed/ratio for JSON events
		Compression: kafka.Snappy,
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
