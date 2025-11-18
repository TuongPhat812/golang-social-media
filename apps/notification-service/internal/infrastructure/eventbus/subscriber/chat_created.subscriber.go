package subscriber

import (
	"context"
	"encoding/json"
	"errors"

	"golang-social-media/apps/notification-service/internal/application/command"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus/subscriber/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/segmentio/kafka-go"
)

var _ contracts.ChatCreatedSubscriber = (*ChatCreatedSubscriber)(nil)

type ChatCreatedSubscriber struct {
	handler *command.HandleChatCreatedCommandHandler
	reader  *kafka.Reader
}

func NewChatCreatedSubscriber(brokers []string, groupID string, handler *command.HandleChatCreatedCommandHandler) (*ChatCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	logger.Component("notification.subscriber.chat_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicChatCreated).
		Msg("chat subscriber configured")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicChatCreated,
	})

	return &ChatCreatedSubscriber{handler: handler, reader: reader}, nil
}

func (s *ChatCreatedSubscriber) Consume(ctx context.Context) {
	logger.Component("notification.subscriber.chat_created").
		Info().
		Str("topic", events.TopicChatCreated).
		Msg("starting chat consumer")
	go func() {
		for {
			msg, err := s.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
					logger.Component("notification.subscriber.chat_created").
						Info().
						Msg("chat consumer shutting down")
					return
				}
				logger.Component("notification.subscriber.chat_created").
					Error().
					Err(err).
					Msg("failed to read ChatCreated message")
				continue
			}

			logger.Component("notification.subscriber.chat_created").
				Info().
				Str("topic", msg.Topic).
				Int("partition", msg.Partition).
				Int64("offset", msg.Offset).
				Msg("received ChatCreated message")

			var event events.ChatCreated
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Component("notification.subscriber.chat_created").
					Error().
					Err(err).
					Msg("failed to unmarshal ChatCreated event")
				continue
			}

			if err := s.handler.Execute(ctx, event); err != nil {
				logger.Component("notification.subscriber.chat_created").
					Error().
					Err(err).
					Msg("failed to handle ChatCreated event")
			} else {
				logger.Component("notification.subscriber.chat_created").
					Info().
					Str("message_id", event.Message.ID).
					Str("sender_id", event.Message.SenderID).
					Str("receiver_id", event.Message.ReceiverID).
					Msg("successfully processed ChatCreated event")
			}
		}
	}()
}

func (s *ChatCreatedSubscriber) Close() error {
	return s.reader.Close()
}
