package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/myself/golang-social-media/common/events"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	log.Printf("[notification-service] initializing kafka publisher with brokers: %v", brokers)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicNotificationCreated,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("[notification-service] publish NotificationCreated to Kafka topic %s: %s", events.TopicNotificationCreated, string(payload))

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.NotificationID),
		Value: payload,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
