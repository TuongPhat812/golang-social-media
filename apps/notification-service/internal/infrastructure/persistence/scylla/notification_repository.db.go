package scylla

import (
	"context"

	"github.com/gocql/gocql"
	"golang-social-media/apps/notification-service/internal/domain/notification"
)

type NotificationRepository struct {
	session *gocql.Session
}

func NewNotificationRepository(session *gocql.Session) *NotificationRepository {
	return &NotificationRepository{session: session}
}

func (r *NotificationRepository) Insert(ctx context.Context, n notification.Notification) error {
	return r.session.Query(`INSERT INTO notifications_by_user (user_id, created_at, notification_id, type, title, body, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.UserID, n.CreatedAt, n.ID, string(n.Type), n.Title, n.Body, n.Metadata,
	).WithContext(ctx).Exec()
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit int) ([]notification.Notification, error) {
	iter := r.session.Query(`SELECT user_id, created_at, notification_id, type, title, body, metadata
		FROM notifications_by_user WHERE user_id = ? LIMIT ?`,
		userID, limit,
	).WithContext(ctx).Iter()

	var results []notification.Notification
	var n notification.Notification
	var typ string
	for iter.Scan(&n.UserID, &n.CreatedAt, &n.ID, &typ, &n.Title, &n.Body, &n.Metadata) {
		n.Type = notification.Type(typ)
		results = append(results, n)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return results, nil
}
