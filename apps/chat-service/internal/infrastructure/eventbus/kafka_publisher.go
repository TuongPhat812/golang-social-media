package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"golang-social-media/pkg/events"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	log.Printf("[chat-service] initializing kafka publisher with brokers: %v", brokers)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicChatCreated,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishChatCreated(ctx context.Context, event events.ChatCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("[chat-service] publish ChatCreated to Kafka topic %s: %s", events.TopicChatCreated, string(payload))

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Message.ID),
		Value: payload,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
