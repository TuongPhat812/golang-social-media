package notifications

import (
	"context"
	"time"

	domainuser "golang-social-media/apps/notification-service/internal/domain/user"
	"golang-social-media/pkg/events"
)

type EventPublisher interface {
	PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error
}

type Service interface {
	SampleUser() domainuser.User
	HandleChatCreated(ctx context.Context, event events.ChatCreated) error
}

type service struct {
	publisher EventPublisher
}

func NewService(publisher EventPublisher) Service {
	return &service{publisher: publisher}
}

func (s *service) SampleUser() domainuser.User {
	return domainuser.User{ID: "2", Username: "notify", FullName: "Notification Service"}
}

func (s *service) HandleChatCreated(ctx context.Context, event events.ChatCreated) error {
	notification := events.NotificationCreated{
		NotificationID: "noti-" + time.Now().UTC().Format("20060102150405"),
		Recipient: events.NotificationRecipient{ID: event.Message.ReceiverID},
		Message:        "New chat message from " + event.Message.SenderID,
		CreatedAt:      time.Now().UTC(),
	}

	return s.publisher.PublishNotificationCreated(ctx, notification)
}
