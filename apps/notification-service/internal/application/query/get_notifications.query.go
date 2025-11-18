package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/notification-service/internal/application/query/contracts"
	"golang-social-media/apps/notification-service/internal/domain/notification"
	"golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/logger"
)

var _ contracts.GetNotificationsQuery = (*getNotificationsQuery)(nil)

type getNotificationsQuery struct {
	repo *scylla.NotificationRepository
	log  *zerolog.Logger
}

func NewGetNotificationsQuery(repo *scylla.NotificationRepository) contracts.GetNotificationsQuery {
	return &getNotificationsQuery{
		repo: repo,
		log:  logger.Component("notification.query.get_notifications"),
	}
}

func (q *getNotificationsQuery) Execute(ctx context.Context, userID string, limit int) ([]notification.Notification, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	notifications, err := q.repo.ListByUser(ctx, userID, limit)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Int("limit", limit).
			Msg("failed to get notifications")
		return nil, err
	}

	q.log.Info().
		Str("user_id", userID).
		Int("count", len(notifications)).
		Msg("notifications retrieved")

	return notifications, nil
}

