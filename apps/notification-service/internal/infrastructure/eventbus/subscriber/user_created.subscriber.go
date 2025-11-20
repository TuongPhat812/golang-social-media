package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/notification-service/internal/application/command"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus/subscriber/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.UserCreatedSubscriber = (*UserCreatedSubscriber)(nil)

type UserCreatedSubscriber struct {
	reader  *kafka.Reader
	handler *command.HandleUserCreatedCommandHandler
}

func NewUserCreatedSubscriber(brokers []string, groupID string, handler *command.HandleUserCreatedCommandHandler) (*UserCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("notification.subscriber.user_created").
		Info().
		Strs("brokers_before_reader", brokers).
		Msg("creating kafka reader with brokers")
	
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    events.TopicUserCreated,
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

	logger.Component("notification.subscriber.user_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicUserCreated).
		Msg("user subscriber configured")

	return &UserCreatedSubscriber{
		reader:  reader,
		handler: handler,
	}, nil
}

func (s *UserCreatedSubscriber) Consume(ctx context.Context) {
	logger.Component("notification.subscriber.user_created").
		Info().
		Str("topic", events.TopicUserCreated).
		Msg("starting user consumer")
	go func() {
		for {
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
					logger.Component("notification.subscriber.user_created").
						Info().
						Msg("user consumer shutting down")
					return
				}
				logger.Component("notification.subscriber.user_created").
					Error().
					Err(err).
					Msg("failed to read UserCreated message")
				continue
			}

			logger.Component("notification.subscriber.user_created").
				Info().
				Str("topic", msg.Topic).
				Int("partition", msg.Partition).
				Int64("offset", msg.Offset).
				Msg("received UserCreated message")

			var event events.UserCreated
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("notification.subscriber.user_created").
					Error().
					Err(err).
					Msg("failed to decode UserCreated event")
				continue
			}

			if err := s.handler.Execute(ctx, event); err != nil {
				logger.Component("notification.subscriber.user_created").
					Error().
					Err(err).
					Msg("failed to handle UserCreated event")
			} else {
				logger.Component("notification.subscriber.user_created").
					Info().
					Str("user_id", event.ID).
					Msg("successfully processed UserCreated event")
			}
		}
	}()
}

func (s *UserCreatedSubscriber) Close() error {
	return s.reader.Close()
}

