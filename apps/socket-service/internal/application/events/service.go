package events

import (
	"context"

	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type Broadcaster interface {
	BroadcastChatCreated(event events.ChatCreated)
	BroadcastNotificationCreated(event events.NotificationCreated)
}

type Service interface {
	HandleChatCreated(ctx context.Context, event events.ChatCreated) error
	HandleNotificationCreated(ctx context.Context, event events.NotificationCreated) error
}

type service struct {
	broadcaster Broadcaster
}

func NewService(broadcaster Broadcaster) Service {
	return &service{broadcaster: broadcaster}
}

func (s *service) HandleChatCreated(ctx context.Context, event events.ChatCreated) error {
	logger.Info().
		Str("topic", events.TopicChatCreated).
		Str("message_id", event.Message.ID).
		Msg("socket-service received ChatCreated event")
	s.broadcaster.BroadcastChatCreated(event)
	return nil
}

func (s *service) HandleNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	logger.Info().
		Str("topic", events.TopicNotificationCreated).
		Str("notification_id", event.NotificationID).
		Msg("socket-service received NotificationCreated event")
	s.broadcaster.BroadcastNotificationCreated(event)
	return nil
}
