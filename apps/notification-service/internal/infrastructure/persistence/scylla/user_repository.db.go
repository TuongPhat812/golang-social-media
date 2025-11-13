package scylla

import (
	"context"

	"github.com/gocql/gocql"
	domainuser "golang-social-media/apps/notification-service/internal/domain/user"
)

type UserRepository struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{session: session}
}

func (r *UserRepository) Upsert(ctx context.Context, user domainuser.User) error {
	return r.session.Query(`INSERT INTO notification_users (user_id, email, name, created_at)
		VALUES (?, ?, ?, ?)`,
		user.ID, user.Email, user.Name, user.CreatedAt,
	).WithContext(ctx).Exec()
}
