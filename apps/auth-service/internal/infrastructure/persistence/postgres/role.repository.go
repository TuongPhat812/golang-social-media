package postgres

import (
	"errors"

	"golang-social-media/apps/auth-service/internal/domain/role"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres/mappers"
	pkgerrors "golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound      = pkgerrors.NewNotFoundError("role_not_found")
	ErrRoleAlreadyExists = pkgerrors.NewConflictError("role_already_exists")
)

type RoleRepository struct {
	db     *gorm.DB
	mapper *mappers.RoleMapper
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db:     db,
		mapper: mappers.NewRoleMapper(),
	}
}

func (r *RoleRepository) Create(roleEntity role.Role) error {
	model := r.mapper.FromDomain(roleEntity)
	if err := r.db.Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrRoleAlreadyExists
		}
		logger.Component("auth.persistence.role_repository").
			Error().
			Err(err).
			Str("role_id", roleEntity.ID).
			Str("name", roleEntity.Name).
			Msg("failed to create role")
		return err
	}
	return nil
}

func (r *RoleRepository) GetByID(id string) (role.Role, error) {
	var model RoleModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role.Role{}, ErrRoleNotFound
		}
		logger.Component("auth.persistence.role_repository").
			Error().
			Err(err).
			Str("role_id", id).
			Msg("failed to get role by ID")
		return role.Role{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *RoleRepository) GetByName(name string) (role.Role, error) {
	var model RoleModel
	if err := r.db.Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role.Role{}, ErrRoleNotFound
		}
		logger.Component("auth.persistence.role_repository").
			Error().
			Err(err).
			Str("name", name).
			Msg("failed to get role by name")
		return role.Role{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *RoleRepository) Update(roleEntity role.Role) error {
	model := r.mapper.FromDomain(roleEntity)
	if err := r.db.Model(&RoleModel{}).Where("id = ?", roleEntity.ID).Updates(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrRoleAlreadyExists
		}
		logger.Component("auth.persistence.role_repository").
			Error().
			Err(err).
			Str("role_id", roleEntity.ID).
			Msg("failed to update role")
		return err
	}
	return nil
}

func (r *RoleRepository) List() ([]role.Role, error) {
	var models []RoleModel
	if err := r.db.Find(&models).Error; err != nil {
		logger.Component("auth.persistence.role_repository").
			Error().
			Err(err).
			Msg("failed to list roles")
		return nil, err
	}

	roles := make([]role.Role, len(models))
	for i, model := range models {
		roles[i] = r.mapper.ToDomain(model)
	}
	return roles, nil
}

