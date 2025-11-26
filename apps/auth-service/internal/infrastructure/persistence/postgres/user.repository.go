package postgres

import (
	"context"
	"errors"

	"golang-social-media/apps/auth-service/internal/domain/user"
	authcache "golang-social-media/apps/auth-service/internal/infrastructure/cache"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres/mappers"
	pkgcache "golang-social-media/pkg/cache"
	pkgerrors "golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrEmailExists  = pkgerrors.NewConflictError(pkgerrors.CodeEmailAlreadyExists)
	ErrInvalidAuth  = pkgerrors.NewValidationError(pkgerrors.CodeInvalidCredentials, nil)
	ErrUserNotFound = pkgerrors.NewNotFoundError(pkgerrors.CodeUserNotFound)
)

type UserRepository struct {
	db     *gorm.DB
	mapper *mappers.UserMapper
	cache  *authcache.UserCache
}

func NewUserRepository(db *gorm.DB, cache *authcache.UserCache) *UserRepository {
	return &UserRepository{
		db:     db,
		mapper: mappers.NewUserMapper(),
		cache:  cache,
	}
}

// NewUserRepositoryWithTx creates a repository that uses a specific transaction
// This is used within UnitOfWork to ensure all operations share the same transaction
func NewUserRepositoryWithTx(tx *gorm.DB, mapper mappers.UserMapper, cache *authcache.UserCache) *UserRepository {
	return &UserRepository{
		db:     tx,
		mapper: &mapper,
		cache:  cache, // Cache is typically not used in transactions
	}
}

func (r *UserRepository) Create(u user.User) error {
	// Check cache first (optimistic check)
	if r.cache != nil {
		if _, err := r.cache.GetUserByEmail(context.Background(), u.Email); err == nil {
			return ErrEmailExists
		}
	}

	model := r.mapper.FromDomain(u)
	if err := r.db.Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrEmailExists
		}
		logger.Component("auth.persistence.user_repository").
			Error().
			Err(err).
			Str("user_id", u.ID).
			Str("email", u.Email).
			Msg("failed to create user")
		return err
	}

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

	var model UserModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, ErrInvalidAuth
		}
		logger.Component("auth.persistence.user_repository").
			Error().
			Err(err).
			Str("email", email).
			Msg("failed to find user by email")
		return user.User{}, err
	}

	u := r.mapper.ToDomain(model)

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

	var model UserModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, ErrUserNotFound
		}
		logger.Component("auth.persistence.user_repository").
			Error().
			Err(err).
			Str("user_id", id).
			Msg("failed to get user by ID")
		return user.User{}, err
	}

	u := r.mapper.ToDomain(model)

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
	// Invalidate cache before update
	if r.cache != nil {
		// Get old user to invalidate cache
		var oldModel UserModel
		if err := r.db.Where("id = ?", u.ID).First(&oldModel).Error; err == nil {
			oldUser := r.mapper.ToDomain(oldModel)
			if err := r.cache.DeleteUser(context.Background(), oldUser.ID, oldUser.Email); err != nil {
				logger.Component("auth.persistence.user_repository").
					Warn().
					Err(err).
					Str("user_id", u.ID).
					Msg("failed to delete user from cache before update")
			}
		}
	}

	model := r.mapper.FromDomain(u)
	if err := r.db.Model(&UserModel{}).Where("id = ?", u.ID).Updates(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		logger.Component("auth.persistence.user_repository").
			Error().
			Err(err).
			Str("user_id", u.ID).
			Msg("failed to update user")
		return err
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

