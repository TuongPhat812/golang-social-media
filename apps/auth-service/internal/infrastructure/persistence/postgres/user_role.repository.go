package postgres

import (
	"golang-social-media/apps/auth-service/internal/domain/user_role"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}

func (r *UserRoleRepository) Create(userRole user_role.UserRole) error {
	model := UserRoleModel{
		UserID:    userRole.UserID,
		RoleID:    userRole.RoleID,
		CreatedAt: userRole.CreatedAt,
	}
	// Use FirstOrCreate to handle duplicate key gracefully
	if err := r.db.Where("user_id = ? AND role_id = ?", userRole.UserID, userRole.RoleID).
		FirstOrCreate(&model).Error; err != nil {
		logger.Component("auth.persistence.user_role_repository").
			Error().
			Err(err).
			Str("user_id", userRole.UserID).
			Str("role_id", userRole.RoleID).
			Msg("failed to create user role")
		return err
	}
	return nil
}

func (r *UserRoleRepository) GetUserRoles(userID string) ([]string, error) {
	var models []UserRoleModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		logger.Component("auth.persistence.user_role_repository").
			Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user roles")
		return nil, err
	}

	roleIDs := make([]string, len(models))
	for i, model := range models {
		roleIDs[i] = model.RoleID
	}
	return roleIDs, nil
}

func (r *UserRoleRepository) GetRoleUsers(roleID string) ([]string, error) {
	var models []UserRoleModel
	if err := r.db.Where("role_id = ?", roleID).Find(&models).Error; err != nil {
		logger.Component("auth.persistence.user_role_repository").
			Error().
			Err(err).
			Str("role_id", roleID).
			Msg("failed to get role users")
		return nil, err
	}

	userIDs := make([]string, len(models))
	for i, model := range models {
		userIDs[i] = model.UserID
	}
	return userIDs, nil
}

func (r *UserRoleRepository) Delete(userID, roleID string) error {
	if err := r.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&UserRoleModel{}).Error; err != nil {
		logger.Component("auth.persistence.user_role_repository").
			Error().
			Err(err).
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("failed to delete user role")
		return err
	}
	return nil
}

func (r *UserRoleRepository) HasRole(userID, roleID string) (bool, error) {
	var count int64
	if err := r.db.Model(&UserRoleModel{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		logger.Component("auth.persistence.user_role_repository").
			Error().
			Err(err).
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("failed to check user role")
		return false, err
	}
	return count > 0, nil
}

