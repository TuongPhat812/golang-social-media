package eventbus

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/notification-service/internal/application/consumers"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type UserSubscriber struct {
	reader  *kafka.Reader
	handler *consumers.UserCreatedConsumer
}

func NewUserSubscriber(brokers []string, groupID string, handler *consumers.UserCreatedConsumer) (*UserSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicUserCreated,
	})

	logger.Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicUserCreated).
		Msg("notification-service user subscriber configured")

	return &UserSubscriber{
		reader:  reader,
		handler: handler,
	}, nil
}

func (s *UserSubscriber) Consume(ctx context.Context) {
	go func() {
		for {
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
					logger.Info().Msg("notification-service user consumer shutting down")
					return
				}
				logger.Error().Err(err).Msg("notification-service failed to read UserCreated message")
				continue
			}

			var event events.UserCreated
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error().Err(err).Msg("notification-service failed to decode UserCreated event")
				continue
			}

			if err := s.handler.Handle(ctx, event); err != nil {
				logger.Error().Err(err).Msg("notification-service failed to handle UserCreated event")
			}
		}
	}()
}

func (s *UserSubscriber) Close() error {
	return s.reader.Close()
}
