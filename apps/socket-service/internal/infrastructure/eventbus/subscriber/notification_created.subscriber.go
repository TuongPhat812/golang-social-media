package subscriber

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/socket-service/internal/infrastructure/eventbus/subscriber/contracts"
	appevents "golang-social-media/apps/socket-service/internal/application/events"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.NotificationCreatedSubscriber = (*NotificationCreatedSubscriber)(nil)

type NotificationCreatedSubscriber struct {
	reader       *kafka.Reader
	eventHandler appevents.Service
	log          *zerolog.Logger
}

func NewNotificationCreatedSubscriber(
	brokers []string,
	groupID string,
	eventHandler appevents.Service,
) (*NotificationCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicNotificationCreated,
	})

	logger.Component("socket.subscriber.notification_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicNotificationCreated).
		Msg("notification subscriber configured")

	return &NotificationCreatedSubscriber{
		reader:       reader,
		eventHandler: eventHandler,
		log:          logger.Component("socket.subscriber.notification_created"),
	}, nil
}

func (s *NotificationCreatedSubscriber) Consume(ctx context.Context) {
	s.log.Info().
		Str("topic", events.TopicNotificationCreated).
		Msg("starting notification consumer")

	for {
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				s.log.Info().Msg("notification listener shutting down")
				return
			}
			s.log.Error().
				Err(err).
				Msg("notification listener error")
			continue
		}

		var event events.NotificationCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			s.log.Error().
				Err(err).
				Msg("failed to decode NotificationCreated event")
			continue
		}

		s.log.Info().
			Str("notification_id", event.Notification.ID).
			Str("user_id", event.Notification.UserID).
			Msg("received NotificationCreated message")

		if err := s.eventHandler.HandleNotificationCreated(ctx, event); err != nil {
			s.log.Error().
				Err(err).
				Str("notification_id", event.Notification.ID).
				Msg("failed to handle NotificationCreated event")
		} else {
			s.log.Info().
				Str("notification_id", event.Notification.ID).
				Msg("successfully processed NotificationCreated event")
		}
	}
}

func (s *NotificationCreatedSubscriber) Close() error {
	return s.reader.Close()
}
