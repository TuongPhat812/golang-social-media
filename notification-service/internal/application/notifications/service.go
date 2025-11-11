package notifications

import (
	"context"
	"time"

	"github.com/myself/golang-social-media/common/domain/user"
	"github.com/myself/golang-social-media/common/events"
)

type EventPublisher interface {
	PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error
}

type Service interface {
	SampleUser() user.User
	HandleChatCreated(ctx context.Context, event events.ChatCreated) error
}

type service struct {
	publisher EventPublisher
}

func NewService(publisher EventPublisher) Service {
	return &service{publisher: publisher}
}

func (s *service) SampleUser() user.User {
	return user.User{ID: "2", Username: "notify", FullName: "Notification Service"}
}

func (s *service) HandleChatCreated(ctx context.Context, event events.ChatCreated) error {
	notification := events.NotificationCreated{
		NotificationID: "noti-" + time.Now().UTC().Format("20060102150405"),
		Recipient:      user.User{ID: event.Message.ReceiverID},
		Message:        "New chat message from " + event.Message.SenderID,
		CreatedAt:      time.Now().UTC(),
	}

	return s.publisher.PublishNotificationCreated(ctx, notification)
}
