package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	"golang-social-media/apps/notification-service/internal/application/command/dto"
	domainnotification "golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.HandleChatCreatedCommand = (*HandleChatCreatedCommandHandler)(nil)

type HandleChatCreatedCommandHandler struct {
	createNotificationCmd contracts.CreateNotificationCommand
	log                   *zerolog.Logger
}

func NewHandleChatCreatedCommand(createNotificationCmd contracts.CreateNotificationCommand) *HandleChatCreatedCommandHandler {
	return &HandleChatCreatedCommandHandler{
		createNotificationCmd: createNotificationCmd,
		log:                   logger.Component("notification.command.handle_chat_created"),
	}
}

func (c *HandleChatCreatedCommandHandler) Execute(ctx context.Context, event events.ChatCreated) error {
	_, err := c.createNotificationCmd.Execute(ctx, dto.CreateNotificationCommandRequest{
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

