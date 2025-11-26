package repository

import (
	"golang-social-media/apps/auth-service/internal/domain/user"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Create(u user.User) error
	GetByID(id string) (user.User, error)
	FindByEmail(email string) (user.User, error)
	Update(u user.User) error
}

