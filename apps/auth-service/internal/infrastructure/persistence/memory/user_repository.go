package memory

import (
	"errors"
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/user"
)

var (
	ErrEmailExists   = errors.New("email already registered")
	ErrInvalidAuth   = errors.New("invalid credentials")
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("token not found")
)

type UserRepository struct {
	mu       sync.RWMutex
	byID     map[string]user.User
	byEmail  map[string]user.User
	password map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:     make(map[string]user.User),
		byEmail:  make(map[string]user.User),
		password: make(map[string]string),
	}
}

func (r *UserRepository) Create(u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[u.Email]; exists {
		return ErrEmailExists
	}

	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	r.password[u.ID] = u.Password
	return nil
}

func (r *UserRepository) FindByEmail(email string) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byEmail[email]
	if !ok {
		return user.User{}, ErrInvalidAuth
	}
	return u, nil
}

func (r *UserRepository) GetByID(id string) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byID[id]
	if !ok {
		return user.User{}, ErrUserNotFound
	}
	return u, nil
}
