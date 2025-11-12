package eventbus

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/notification-service/internal/application/notifications"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type Subscriber struct {
	service notifications.Service
	reader  *kafka.Reader
}

func NewSubscriber(brokers []string, groupID string, service notifications.Service) (*Subscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Msg("notification-service kafka subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicChatCreated,
	})

	return &Subscriber{service: service, reader: reader}, nil
}

func (s *Subscriber) ConsumeChatCreated(ctx context.Context) {
	go func() {
		for {
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
					logger.Info().Msg("notification-service chat consumer shutting down")
					return
				}
				logger.Error().Err(err).Msg("notification-service failed to read ChatCreated message")
				continue
			}

			var event events.ChatCreated
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error().Err(err).Msg("notification-service failed to unmarshal ChatCreated event")
				continue
			}

			if err := s.service.HandleChatCreated(ctx, event); err != nil {
				logger.Error().Err(err).Msg("notification-service failed to handle ChatCreated event")
			}
		}
	}()
}

func (s *Subscriber) Close() error {
	return s.reader.Close()
}
