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

var _ contracts.ChatCreatedSubscriber = (*ChatCreatedSubscriber)(nil)

type ChatCreatedSubscriber struct {
	reader      *kafka.Reader
	eventHandler appevents.Service
	log         *zerolog.Logger
}

func NewChatCreatedSubscriber(
	brokers []string,
	groupID string,
	eventHandler appevents.Service,
) (*ChatCreatedSubscriber, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}
	if groupID == "" {
		return nil, errors.New("groupID must be provided")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   events.TopicChatCreated,
	})

	logger.Component("socket.subscriber.chat_created").
		Info().
		Strs("brokers", brokers).
		Str("group", groupID).
		Str("topic", events.TopicChatCreated).
		Msg("chat subscriber configured")

	return &ChatCreatedSubscriber{
		reader:       reader,
		eventHandler: eventHandler,
		log:          logger.Component("socket.subscriber.chat_created"),
	}, nil
}

func (s *ChatCreatedSubscriber) Consume(ctx context.Context) {
	s.log.Info().
		Str("topic", events.TopicChatCreated).
		Msg("starting chat consumer")

	for {
		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, kafka.ErrGroupClosed) {
				s.log.Info().Msg("chat listener shutting down")
				return
			}
			s.log.Error().
				Err(err).
				Msg("chat listener error")
			continue
		}

		var event events.ChatCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			s.log.Error().
				Err(err).
				Msg("failed to decode ChatCreated event")
			continue
		}

		s.log.Info().
			Str("message_id", event.Message.ID).
			Str("sender_id", event.Message.SenderID).
			Str("receiver_id", event.Message.ReceiverID).
			Msg("received ChatCreated message")

		if err := s.eventHandler.HandleChatCreated(ctx, event); err != nil {
			s.log.Error().
				Err(err).
				Str("message_id", event.Message.ID).
				Msg("failed to handle ChatCreated event")
		} else {
			s.log.Info().
				Str("message_id", event.Message.ID).
				Msg("successfully processed ChatCreated event")
		}
	}
}

func (s *ChatCreatedSubscriber) Close() error {
	return s.reader.Close()
}
