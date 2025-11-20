package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/socket-service/internal/infrastructure/eventbus/subscriber/contracts"
	appevents "golang-social-media/apps/socket-service/internal/application/events"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.ChatCreatedSubscriber = (*ChatCreatedSubscriber)(nil)

type ChatCreatedSubscriber struct {
	reader      *kafka.Reader
	eventHandler appevents.Service
	log         *zerolog.Logger
}

func NewChatCreatedSubscriber(
	brokers []string,
	groupID string,
	eventHandler appevents.Service,
) (*ChatCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    events.TopicChatCreated,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// Connection timeouts
		Dialer: &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			KeepAlive:     5 * time.Minute,
		},
		// Read timeouts
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 1 * time.Second,
		// Commit interval - commit offsets every 1 second
		CommitInterval: 1 * time.Second,
	})

	logger.Component("socket.subscriber.chat_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicChatCreated).
		Msg("chat subscriber configured")

	return &ChatCreatedSubscriber{
		reader:       reader,
		eventHandler: eventHandler,
		log:          logger.Component("socket.subscriber.chat_created"),
	}, nil
}

func (s *ChatCreatedSubscriber) Consume(ctx context.Context) {
	s.log.Info().
		Str("topic", events.TopicChatCreated).
		Msg("starting chat consumer (this may take a few seconds to connect to Kafka)")

	for {
			msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				s.log.Info().Msg("chat listener shutting down")
				return
			}
			// Log first connection attempt separately
			if err != nil {
				s.log.Error().
					Err(err).
					Msg("chat listener error")
				// Add small delay on error to avoid tight loop
				select {
				case <-ctx.Done():
					return
				case <-time.After(100 * time.Millisecond):
				}
				continue
			}
			
			// Log first successful message read
			if msg.Offset == 0 || msg.Partition == 0 {
				s.log.Info().
					Str("topic", msg.Topic).
					Int("partition", msg.Partition).
					Int64("offset", msg.Offset).
					Msg("chat consumer connected and reading messages")
			}
		}

		var event events.ChatCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			s.log.Error().
				Err(err).
				Msg("failed to decode ChatCreated event")
			continue
		}

		s.log.Info().
			Str("message_id", event.Message.ID).
			Str("sender_id", event.Message.SenderID).
			Str("receiver_id", event.Message.ReceiverID).
			Msg("received ChatCreated message")

		if err := s.eventHandler.HandleChatCreated(ctx, event); err != nil {
			s.log.Error().
				Err(err).
				Str("message_id", event.Message.ID).
				Msg("failed to handle ChatCreated event")
		} else {
			s.log.Info().
				Str("message_id", event.Message.ID).
				Msg("successfully processed ChatCreated event")
		}
	}
}

func (s *ChatCreatedSubscriber) Close() error {
	return s.reader.Close()
}
