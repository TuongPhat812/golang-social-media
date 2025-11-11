package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/myself/golang-social-media/pkg/events"
	appevents "github.com/myself/golang-social-media/apps/socket-service/internal/application/events"
	"github.com/segmentio/kafka-go"
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

	log.Printf("[socket-service] initializing kafka listeners with brokers: %v, chat group: %s, notification group: %s", brokers, chatGroupID, notificationGroupID)

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

	return &Listener{
		service:            service,
		chatReader:         chatReader,
		notificationReader: notificationReader,
	}, nil
}

func (l *Listener) Start(ctx context.Context) {
	log.Println("[socket-service] starting Kafka listeners")
	go l.consumeChatCreated(ctx)
	go l.consumeNotificationCreated(ctx)
}

func (l *Listener) consumeChatCreated(ctx context.Context) {
	for {
		msg, err := l.chatReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				log.Println("[socket-service] chat listener shutting down")
				return
			}
			log.Printf("[socket-service] chat listener error: %v", err)
			continue
		}

		var event events.ChatCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[socket-service] failed to decode ChatCreated event: %v", err)
			continue
		}

		if err := l.service.HandleChatCreated(ctx, event); err != nil {
			log.Printf("[socket-service] failed to handle ChatCreated event: %v", err)
		}
	}
}

func (l *Listener) consumeNotificationCreated(ctx context.Context) {
	for {
		msg, err := l.notificationReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				log.Println("[socket-service] notification listener shutting down")
				return
			}
			log.Printf("[socket-service] notification listener error: %v", err)
			continue
		}

		var event events.NotificationCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[socket-service] failed to decode NotificationCreated event: %v", err)
			continue
		}

		if err := l.service.HandleNotificationCreated(ctx, event); err != nil {
			log.Printf("[socket-service] failed to handle NotificationCreated event: %v", err)
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
