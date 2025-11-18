package scylla

import (
	"context"
	"time"

	"golang-social-media/apps/notification-service/internal/domain/notification"

	"github.com/gocql/gocql"
)

type NotificationRepository struct {
	session *gocql.Session
}

func NewNotificationRepository(session *gocql.Session) *NotificationRepository {
	return &NotificationRepository{session: session}
}

func (r *NotificationRepository) Insert(ctx context.Context, n notification.Notification) error {
	return r.session.Query(`INSERT INTO notifications_by_user (user_id, created_at, notification_id, type, title, body, metadata, read_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID, n.CreatedAt, n.ID, string(n.Type), n.Title, n.Body, n.Metadata, n.ReadAt,
	).WithContext(ctx).Exec()
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit int) ([]notification.Notification, error) {
	iter := r.session.Query(`SELECT user_id, created_at, notification_id, type, title, body, metadata, read_at
		FROM notifications_by_user WHERE user_id = ? LIMIT ?`,
		userID, limit,
	).WithContext(ctx).Iter()

	var results []notification.Notification
	var n notification.Notification
	var typ string
	var readAt *time.Time
	for iter.Scan(&n.UserID, &n.CreatedAt, &n.ID, &typ, &n.Title, &n.Body, &n.Metadata, &readAt) {
		n.Type = notification.Type(typ)
		n.ReadAt = readAt
		results = append(results, n)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *NotificationRepository) FindByID(ctx context.Context, userID string, notificationID gocql.UUID) (notification.Notification, error) {
	var n notification.Notification
	var typ string
	var readAt *time.Time

	err := r.session.Query(`SELECT user_id, created_at, notification_id, type, title, body, metadata, read_at
		FROM notifications_by_user WHERE user_id = ? AND notification_id = ?`,
		userID, notificationID,
	).WithContext(ctx).Scan(&n.UserID, &n.CreatedAt, &n.ID, &typ, &n.Title, &n.Body, &n.Metadata, &readAt)

	if err != nil {
		return notification.Notification{}, err
	}

	n.Type = notification.Type(typ)
	n.ReadAt = readAt
	return n, nil
}

func (r *NotificationRepository) UpdateReadAt(ctx context.Context, userID string, notificationID gocql.UUID, readAt time.Time) error {
	return r.session.Query(`UPDATE notifications_by_user SET read_at = ?
		WHERE user_id = ? AND notification_id = ?`,
		readAt, userID, notificationID,
	).WithContext(ctx).Exec()
}
