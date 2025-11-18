package command

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"golang-social-media/apps/notification-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/notification-service/internal/application/event_dispatcher"
	"golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/logger"
)

var _ contracts.MarkNotificationReadCommand = (*markNotificationReadCommand)(nil)

type markNotificationReadCommand struct {
	repo           *scylla.NotificationRepository
	eventDispatcher *event_dispatcher.Dispatcher
	log            *zerolog.Logger
}

func NewMarkNotificationReadCommand(
	repo *scylla.NotificationRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.MarkNotificationReadCommand {
	return &markNotificationReadCommand{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("notification.command.mark_notification_read"),
	}
}

func (c *markNotificationReadCommand) Execute(ctx context.Context, userID string, notificationID string) error {
	// Parse notification ID
	uuid, err := gocql.ParseUUID(notificationID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("notification_id", notificationID).
			Msg("invalid notification ID")
		return err
	}

	// Load notification from repository
	notificationModel, err := c.repo.FindByID(ctx, userID, uuid)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Str("notification_id", notificationID).
			Msg("failed to find notification")
		return err
	}

	// Domain logic: mark as read (this adds domain events internally)
	if err := notificationModel.MarkAsRead(); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Str("notification_id", notificationID).
			Msg("failed to mark notification as read")
		return err
	}

	// Persist to database
	if err := c.repo.UpdateReadAt(ctx, userID, uuid, *notificationModel.ReadAt); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Str("notification_id", notificationID).
			Msg("failed to update notification read_at")
		return err
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
				Str("notification_id", notificationID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("notification_id", notificationID).
		Str("user_id", userID).
		Msg("notification marked as read")

	return nil
}

