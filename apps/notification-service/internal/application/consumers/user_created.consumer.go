package consumers

import (
	"context"
	"fmt"

	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	"golang-social-media/apps/notification-service/internal/application/command/dto"
	domainnotification "golang-social-media/apps/notification-service/internal/domain/notification"
	domainuser "golang-social-media/apps/notification-service/internal/domain/user"
	scyllarepo "golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type UserCreatedConsumer struct {
	userRepo              *scyllarepo.UserRepository
	createNotificationCmd contracts.CreateNotificationCommand
}

func NewUserCreatedConsumer(userRepo *scyllarepo.UserRepository, cmd contracts.CreateNotificationCommand) *UserCreatedConsumer {
	return &UserCreatedConsumer{
		userRepo:              userRepo,
		createNotificationCmd: cmd,
	}
}

func (c *UserCreatedConsumer) Handle(ctx context.Context, event events.UserCreated) error {
	if err := c.userRepo.Upsert(ctx, domainuser.User{
		ID:        event.ID,
		Email:     event.Email,
		Name:      event.Name,
		CreatedAt: event.CreatedAt,
	}); err != nil {
		logger.Component("notification.consumer.user_created").
			Error().
			Err(err).
			Str("user_id", event.ID).
			Msg("failed to replicate user")
		return err
	}

	_, err := c.createNotificationCmd.Handle(ctx, dto.CreateNotificationCommandRequest{
		UserID: event.ID,
		Type:   domainnotification.TypeWelcome,
		Title:  "Chào mừng bạn đến với Golang Social Media",
		Body:   fmt.Sprintf("Xin chào %s! Cảm ơn bạn đã đăng ký.", event.Name),
		Time:   event.CreatedAt,
		Metadata: map[string]string{
			"email": event.Email,
		},
	})

	return err
}
