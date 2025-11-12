package eventbus

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	appevents "golang-social-media/apps/socket-service/internal/application/events"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type Listener struct {
	service            appevents.Service
	chatReader         *kafka.Reader
	notificationReader *kafka.Reader
}

func NewListener(brokers []string, chatGroupID, notificationGroupID string, service appevents.Service) (*Listener, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if chatGroupID == "" {
		return nil, errors.New("chatGroupID must be provided")
	}
	if notificationGroupID == "" {
		return nil, errors.New("notificationGroupID must be provided")
	}

	chatReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: chatGroupID,
		Topic:   events.TopicChatCreated,
	})

	notificationReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: notificationGroupID,
		Topic:   events.TopicNotificationCreated,
	})

	logger.Info().
		Strs("brokers", brokers).
		Str("chat_group", chatGroupID).
		Str("notification_group", notificationGroupID).
		Msg("socket-service kafka listeners configured")

	return &Listener{
		service:            service,
		chatReader:         chatReader,
		notificationReader: notificationReader,
	}, nil
}

func (l *Listener) Start(ctx context.Context) {
	logger.Info().Msg("socket-service starting Kafka listeners")
	go l.consumeChatCreated(ctx)
	go l.consumeNotificationCreated(ctx)
}

func (l *Listener) consumeChatCreated(ctx context.Context) {
	for {
		msg, err := l.chatReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				logger.Info().Msg("socket-service chat listener shutting down")
				return
			}
			logger.Error().Err(err).Msg("socket-service chat listener error")
			continue
		}

		var event events.ChatCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error().Err(err).Msg("socket-service failed to decode ChatCreated event")
			continue
		}

		if err := l.service.HandleChatCreated(ctx, event); err != nil {
			logger.Error().Err(err).Msg("socket-service failed to handle ChatCreated event")
		}
	}
}

func (l *Listener) consumeNotificationCreated(ctx context.Context) {
	for {
		msg, err := l.notificationReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				logger.Info().Msg("socket-service notification listener shutting down")
				return
			}
			logger.Error().Err(err).Msg("socket-service notification listener error")
			continue
		}

		var event events.NotificationCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error().Err(err).Msg("socket-service failed to decode NotificationCreated event")
			continue
		}

		if err := l.service.HandleNotificationCreated(ctx, event); err != nil {
			logger.Error().Err(err).Msg("socket-service failed to handle NotificationCreated event")
		}
	}
}

func (l *Listener) Close() error {
	var err error
	if l.chatReader != nil {
		if closeErr := l.chatReader.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if l.notificationReader != nil {
		if closeErr := l.notificationReader.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}
