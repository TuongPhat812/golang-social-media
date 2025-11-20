package factories

import "golang-social-media/apps/auth-service/internal/domain/user"

// UserFactory defines the contract for creating User entities
type UserFactory interface {
	CreateUser(email, password, name string) (*user.User, error)
}


