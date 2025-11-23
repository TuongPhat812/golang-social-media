package memory

import (
	"context"
	"errors"
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/user"
	authcache "golang-social-media/apps/auth-service/internal/infrastructure/cache"
	pkgcache "golang-social-media/pkg/cache"
	pkgerrors "golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var (
	ErrEmailExists   = pkgerrors.NewConflictError(pkgerrors.CodeEmailAlreadyExists)
	ErrInvalidAuth   = pkgerrors.NewValidationError(pkgerrors.CodeInvalidCredentials, nil)
	ErrUserNotFound  = pkgerrors.NewNotFoundError(pkgerrors.CodeUserNotFound)
	ErrTokenNotFound = pkgerrors.NewNotFoundError(pkgerrors.CodeTokenInvalid)
)

type UserRepository struct {
	mu       sync.RWMutex
	byID     map[string]user.User
	byEmail  map[string]user.User
	password map[string]string
	cache    *authcache.UserCache
}

func NewUserRepository(cache *authcache.UserCache) *UserRepository {
	return &UserRepository{
		byID:     make(map[string]user.User),
		byEmail:  make(map[string]user.User),
		password: make(map[string]string),
		cache:    cache,
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

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetUser(context.Background(), &u); err != nil {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", u.ID).
				Msg("failed to set user to cache after create")
		}
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (user.User, error) {
	// Check cache first
	if r.cache != nil {
		if u, err := r.cache.GetUserByEmail(context.Background(), email); err == nil {
			logger.Component("auth.persistence.user_repository").
				Debug().
				Str("email", email).
				Msg("user found in cache by email")
			return *u, nil
		} else if !errors.Is(err, pkgcache.ErrCacheMiss) {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("email", email).
				Msg("failed to get user from cache by email")
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byEmail[email]
	if !ok {
		return user.User{}, ErrInvalidAuth
	}

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetUser(context.Background(), &u); err != nil {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", u.ID).
				Msg("failed to set user to cache after db fetch by email")
		}
	}
	return u, nil
}

func (r *UserRepository) GetByID(id string) (user.User, error) {
	// Check cache first
	if r.cache != nil {
		if u, err := r.cache.GetUserByID(context.Background(), id); err == nil {
			logger.Component("auth.persistence.user_repository").
				Debug().
				Str("user_id", id).
				Msg("user found in cache by ID")
			return *u, nil
		} else if !errors.Is(err, pkgcache.ErrCacheMiss) {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", id).
				Msg("failed to get user from cache by ID")
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byID[id]
	if !ok {
		return user.User{}, ErrUserNotFound
	}

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetUser(context.Background(), &u); err != nil {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", u.ID).
				Msg("failed to set user to cache after db fetch by ID")
		}
	}
	return u, nil
}

func (r *UserRepository) Update(u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldUser, exists := r.byID[u.ID]
	if !exists {
		return ErrUserNotFound
	}

	// Invalidate cache before update
	if r.cache != nil {
		// Delete old cache entries
		if err := r.cache.DeleteUser(context.Background(), oldUser.ID, oldUser.Email); err != nil {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", u.ID).
				Msg("failed to delete user from cache before update")
		}
	}

	// Update in-memory storage
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	if u.Password != "" {
		r.password[u.ID] = u.Password
	}

	// Update cache with new data
	if r.cache != nil {
		if err := r.cache.SetUser(context.Background(), &u); err != nil {
			logger.Component("auth.persistence.user_repository").
				Warn().
				Err(err).
				Str("user_id", u.ID).
				Msg("failed to set user to cache after update")
		}
	}

	return nil
}
