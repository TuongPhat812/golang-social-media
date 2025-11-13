package notifications

import (
	"context"
	"time"

	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	"golang-social-media/apps/notification-service/internal/application/command/dto"
	domainnotification "golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/pkg/events"
)

type Service interface {
	HandleChatCreated(ctx context.Context, event events.ChatCreated) error
}

type service struct {
	createNotification contracts.CreateNotificationCommand
}

func NewService(createNotification contracts.CreateNotificationCommand) Service {
	return &service{createNotification: createNotification}
}

func (s *service) HandleChatCreated(ctx context.Context, event events.ChatCreated) error {
	_, err := s.createNotification.Handle(ctx, dto.CreateNotificationCommandRequest{
		UserID: event.Message.ReceiverID,
		Type:   domainnotification.TypeChatMessage,
		Title:  "Tin nhắn mới",
		Body:   "New chat message from " + event.Message.SenderID,
		Time:   time.Now().UTC(),
		Metadata: map[string]string{
			"senderId":   event.Message.SenderID,
			"messageId":  event.Message.ID,
			"content":    event.Message.Content,
			"receiverId": event.Message.ReceiverID,
		},
	})
	return err
}
