package command

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	"golang-social-media/apps/notification-service/internal/application/command/dto"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

type notificationPublisher interface {
	PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error
}

type createNotificationCommand struct {
	repo      *scylla.NotificationRepository
	publisher notificationPublisher
	log       *zerolog.Logger
}

func NewCreateNotificationCommand(
	repo *scylla.NotificationRepository,
	publisher notificationPublisher,
) contracts.CreateNotificationCommand {
	return &createNotificationCommand{
		repo:      repo,
		publisher: publisher,
		log:       logger.Component("notification.command.create_notification"),
	}
}

func (c *createNotificationCommand) Handle(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error) {
	createdAt := req.Time
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	noti := notification.Notification{
		ID:        gocql.TimeUUID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Metadata:  req.Metadata,
		CreatedAt: createdAt,
	}

	if err := c.repo.Insert(ctx, noti); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Str("type", string(req.Type)).
			Msg("failed to persist notification")
		return notification.Notification{}, err
	}

	if c.publisher != nil {
		event := events.NotificationCreated{
			Notification: events.Notification{
				ID:        noti.ID.String(),
				UserID:    noti.UserID,
				Type:      string(noti.Type),
				Title:     noti.Title,
				Body:      noti.Body,
				CreatedAt: noti.CreatedAt,
				Metadata:  noti.Metadata,
			},
		}
		if err := c.publisher.PublishNotificationCreated(ctx, event); err != nil {
			c.log.Error().
				Err(err).
				Str("notification_id", noti.ID.String()).
				Msg("failed to publish NotificationCreated")
			return notification.Notification{}, err
		}
	}

	c.log.Info().
		Str("notification_id", noti.ID.String()).
		Str("user_id", noti.UserID).
		Msg("notification created")

	return noti, nil
}
