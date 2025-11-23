package postgres

import (
	"errors"

	"golang-social-media/apps/auth-service/internal/domain/permission"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres/mappers"
	pkgerrors "golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrPermissionNotFound     = pkgerrors.NewNotFoundError("permission_not_found")
	ErrPermissionAlreadyExists = pkgerrors.NewConflictError("permission_already_exists")
)

type PermissionRepository struct {
	db     *gorm.DB
	mapper *mappers.PermissionMapper
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db:     db,
		mapper: mappers.NewPermissionMapper(),
	}
}

func (r *PermissionRepository) Create(perm permission.Permission) error {
	model := r.mapper.FromDomain(perm)
	if err := r.db.Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrPermissionAlreadyExists
		}
		logger.Component("auth.persistence.permission_repository").
			Error().
			Err(err).
			Str("permission_id", perm.ID).
			Str("resource", perm.Resource).
			Str("action", perm.Action).
			Msg("failed to create permission")
		return err
	}
	return nil
}

func (r *PermissionRepository) GetByID(id string) (permission.Permission, error) {
	var model PermissionModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return permission.Permission{}, ErrPermissionNotFound
		}
		logger.Component("auth.persistence.permission_repository").
			Error().
			Err(err).
			Str("permission_id", id).
			Msg("failed to get permission by ID")
		return permission.Permission{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *PermissionRepository) GetByResourceAction(resource, action string) (permission.Permission, error) {
	var model PermissionModel
	if err := r.db.Where("resource = ? AND action = ?", resource, action).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return permission.Permission{}, ErrPermissionNotFound
		}
		logger.Component("auth.persistence.permission_repository").
			Error().
			Err(err).
			Str("resource", resource).
			Str("action", action).
			Msg("failed to get permission by resource and action")
		return permission.Permission{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *PermissionRepository) List() ([]permission.Permission, error) {
	var models []PermissionModel
	if err := r.db.Find(&models).Error; err != nil {
		logger.Component("auth.persistence.permission_repository").
			Error().
			Err(err).
			Msg("failed to list permissions")
		return nil, err
	}

	permissions := make([]permission.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.mapper.ToDomain(model)
	}
	return permissions, nil
}

