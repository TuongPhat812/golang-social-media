package events

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

// Broadcaster interface for broadcasting events via WebSocket
type Broadcaster interface {
	BroadcastChatCreated(event events.ChatCreated)
	BroadcastNotificationCreated(event events.NotificationCreated)
}

// Service handles events and broadcasts them via WebSocket
type Service interface {
	HandleChatCreated(ctx context.Context, event events.ChatCreated) error
	HandleNotificationCreated(ctx context.Context, event events.NotificationCreated) error
}

type service struct {
	broadcaster Broadcaster
	log         *zerolog.Logger
}

// NewService creates a new event service
func NewService(broadcaster Broadcaster) Service {
	return &service{
		broadcaster: broadcaster,
		log:         logger.Component("socket.events"),
	}
}

func (s *service) HandleChatCreated(ctx context.Context, event events.ChatCreated) error {
	s.log.Info().
		Str("topic", events.TopicChatCreated).
		Str("message_id", event.Message.ID).
		Str("sender_id", event.Message.SenderID).
		Str("receiver_id", event.Message.ReceiverID).
		Msg("handling ChatCreated event")
	s.broadcaster.BroadcastChatCreated(event)
	return nil
}

func (s *service) HandleNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	s.log.Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.Notification.ID).
		Str("user_id", event.Notification.UserID).
		Msg("handling NotificationCreated event")
	s.broadcaster.BroadcastNotificationCreated(event)
	return nil
}
