package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/chat-service/internal/application/command"
	"golang-social-media/apps/chat-service/internal/infrastructure/eventbus/subscriber/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.UserCreatedSubscriber = (*UserCreatedSubscriber)(nil)

type UserCreatedSubscriber struct {
	reader  *kafka.Reader
	handler *command.HandleUserCreatedCommandHandler
}

func NewUserCreatedSubscriber(
	brokers []string,
	groupID string,
	handler *command.HandleUserCreatedCommandHandler,
) (*UserCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("chat.subscriber.user_created").
		Info().
		Strs("brokers", brokers).
		Msg("creating kafka reader for user.created")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicUserCreated,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 5 * time.Minute,
		},
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 1 * time.Second,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})

	logger.Component("chat.subscriber.user_created").
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
	logger.Component("chat.subscriber.user_created").
		Info().
		Str("topic", events.TopicUserCreated).
		Msg("starting user consumer")

	for {
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				logger.Component("chat.subscriber.user_created").
					Info().
					Msg("user consumer shutting down")
				return
			}
			logger.Component("chat.subscriber.user_created").
				Error().
				Err(err).
				Msg("failed to read UserCreated message")
			continue
		}

		logger.Component("chat.subscriber.user_created").
			Info().
			Str("topic", msg.Topic).
			Int("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Msg("received UserCreated message")

		var event events.UserCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Component("chat.subscriber.user_created").
				Error().
				Err(err).
				Msg("failed to decode UserCreated event")
			continue
		}

		if err := s.handler.Execute(ctx, event); err != nil {
			logger.Component("chat.subscriber.user_created").
				Error().
				Err(err).
				Str("user_id", event.ID).
				Msg("failed to handle UserCreated event")
		} else {
			logger.Component("chat.subscriber.user_created").
				Info().
				Str("user_id", event.ID).
				Str("email", event.Email).
				Msg("successfully processed UserCreated event")
		}
	}
}

func (s *UserCreatedSubscriber) Close() error {
	return s.reader.Close()
}

