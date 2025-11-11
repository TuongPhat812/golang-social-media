package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/myself/golang-social-media/common/events"
	"github.com/myself/golang-social-media/notification-service/internal/application/notifications"
	"github.com/segmentio/kafka-go"
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

	log.Printf("[notification-service] initializing kafka subscriber with brokers: %v, group: %s", brokers, groupID)

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
					log.Println("[notification-service] chat consumer shutting down")
					return
				}
				log.Printf("[notification-service] failed to read ChatCreated message: %v", err)
				continue
			}

			var event events.ChatCreated
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("[notification-service] failed to unmarshal ChatCreated event: %v", err)
				continue
			}

			if err := s.service.HandleChatCreated(ctx, event); err != nil {
				log.Printf("[notification-service] failed to handle ChatCreated event: %v", err)
			}
		}
	}()
}

func (s *Subscriber) Close() error {
	return s.reader.Close()
}
