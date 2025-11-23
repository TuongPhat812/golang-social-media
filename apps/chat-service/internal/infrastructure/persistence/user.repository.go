package persistence

import (
	"context"

	chatcache "golang-social-media/apps/chat-service/internal/infrastructure/cache"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	cache *chatcache.UserCache
}

func NewUserRepository(db *gorm.DB, userCache *chatcache.UserCache) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: userCache,
	}
}

// Upsert creates or updates a user
func (r *UserRepository) Upsert(ctx context.Context, user UserModel) error {
	// Save to database
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetUser(ctx, &user); err != nil {
			logger.Component("chat.repository.user").
				Warn().
				Err(err).
				Str("user_id", user.ID).
				Msg("failed to update user cache")
		}
	}

	return nil
}

// FindByID finds a user by ID (with cache)
func (r *UserRepository) FindByID(ctx context.Context, id string) (*UserModel, error) {
	// Try cache first
	if r.cache != nil {
		if user, err := r.cache.GetUser(ctx, id); err == nil {
			return user, nil
		}
	}

	// Cache miss, query database
	var user UserModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetUser(ctx, &user); err != nil {
			logger.Component("chat.repository.user").
				Warn().
				Err(err).
				Str("user_id", id).
				Msg("failed to set user cache")
		}
	}

	return &user, nil
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	// Try cache first
	if r.cache != nil {
		if user, err := r.cache.GetUser(ctx, id); err == nil {
			return user != nil, nil
		}
	}

	// Cache miss, query database
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

