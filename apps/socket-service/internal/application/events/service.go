package events

import (
	"context"
	"log"

	"github.com/myself/golang-social-media/pkg/events"
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
	log.Printf("[socket-service] received ChatCreated event: %+v", event)
	s.broadcaster.BroadcastChatCreated(event)
	return nil
}

func (s *service) HandleNotificationCreated(ctx context.Context, event events.NotificationCreated) error {
	log.Printf("[socket-service] received NotificationCreated event: %+v", event)
	s.broadcaster.BroadcastNotificationCreated(event)
	return nil
}
