package postgres

import (
	"golang-social-media/apps/auth-service/internal/domain/role_permission"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{
		db: db,
	}
}

func (r *RolePermissionRepository) Create(rp role_permission.RolePermission) error {
	model := RolePermissionModel{
		RoleID:       rp.RoleID,
		PermissionID: rp.PermissionID,
		CreatedAt:    rp.CreatedAt,
	}
	// Use FirstOrCreate to handle duplicate key gracefully
	if err := r.db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
		FirstOrCreate(&model).Error; err != nil {
		logger.Component("auth.persistence.role_permission_repository").
			Error().
			Err(err).
			Str("role_id", rp.RoleID).
			Str("permission_id", rp.PermissionID).
			Msg("failed to create role permission")
		return err
	}
	return nil
}

func (r *RolePermissionRepository) GetRolePermissions(roleID string) ([]string, error) {
	var models []RolePermissionModel
	if err := r.db.Where("role_id = ?", roleID).Find(&models).Error; err != nil {
		logger.Component("auth.persistence.role_permission_repository").
			Error().
			Err(err).
			Str("role_id", roleID).
			Msg("failed to get role permissions")
		return nil, err
	}

	permissionIDs := make([]string, len(models))
	for i, model := range models {
		permissionIDs[i] = model.PermissionID
	}
	return permissionIDs, nil
}

func (r *RolePermissionRepository) GetPermissionRoles(permissionID string) ([]string, error) {
	var models []RolePermissionModel
	if err := r.db.Where("permission_id = ?", permissionID).Find(&models).Error; err != nil {
		logger.Component("auth.persistence.role_permission_repository").
			Error().
			Err(err).
			Str("permission_id", permissionID).
			Msg("failed to get permission roles")
		return nil, err
	}

	roleIDs := make([]string, len(models))
	for i, model := range models {
		roleIDs[i] = model.RoleID
	}
	return roleIDs, nil
}

func (r *RolePermissionRepository) Delete(roleID, permissionID string) error {
	if err := r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&RolePermissionModel{}).Error; err != nil {
		logger.Component("auth.persistence.role_permission_repository").
			Error().
			Err(err).
			Str("role_id", roleID).
			Str("permission_id", permissionID).
			Msg("failed to delete role permission")
		return err
	}
	return nil
}

func (r *RolePermissionRepository) HasPermission(roleID, permissionID string) (bool, error) {
	var count int64
	if err := r.db.Model(&RolePermissionModel{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error; err != nil {
		logger.Component("auth.persistence.role_permission_repository").
			Error().
			Err(err).
			Str("role_id", roleID).
			Str("permission_id", permissionID).
			Msg("failed to check role permission")
		return false, err
	}
	return count > 0, nil
}

