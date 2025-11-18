package command

import (
	"context"
	"time"

	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	"golang-social-media/apps/notification-service/internal/application/command/dto"
	event_dispatcher "golang-social-media/apps/notification-service/internal/application/event_dispatcher"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/logger"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
)

var _ contracts.CreateNotificationCommand = (*createNotificationCommand)(nil)

type createNotificationCommand struct {
	repo            *scylla.NotificationRepository
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewCreateNotificationCommand(
	repo *scylla.NotificationRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CreateNotificationCommand {
	return &createNotificationCommand{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("notification.command.create_notification"),
	}
}

func (c *createNotificationCommand) Execute(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error) {
	return c.Handle(ctx, req)
}

func (c *createNotificationCommand) Handle(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error) {
	createdAt := req.Time
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	notificationModel := notification.Notification{
		ID:        gocql.TimeUUID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Metadata:  req.Metadata,
		CreatedAt: createdAt,
	}

	// Validate business rules before persisting or publishing
	if err := notificationModel.Validate(); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Str("type", string(req.Type)).
			Msg("notification validation failed")
		return notification.Notification{}, err
	}

	// Domain logic: create notification (this adds domain events internally)
	notificationModel.Create()

	// Persist to database
	if err := c.repo.Insert(ctx, notificationModel); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Str("type", string(req.Type)).
			Msg("failed to persist notification")
		return notification.Notification{}, err
	}

	// Dispatch domain events AFTER successful persistence
	domainEvents := notificationModel.Events()
	notificationModel.ClearEvents() // Clear events after dispatch

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			// Log error but don't fail the command
			// Events can be retried via outbox pattern in production
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("notification_id", notificationModel.ID.String()).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("notification_id", notificationModel.ID.String()).
		Str("user_id", notificationModel.UserID).
		Msg("notification created")

	return notificationModel, nil
}
