package eventbus

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"golang-social-media/pkg/events"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicUserCreated,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaPublisher{writer: writer}, nil
}

func (p *KafkaPublisher) PublishUserCreated(ctx context.Context, event events.UserCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: payload,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
