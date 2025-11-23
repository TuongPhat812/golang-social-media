package mappers

import (
	"golang-social-media/apps/auth-service/internal/domain/permission"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
)

// PermissionMapper maps between domain Permission and PostgreSQL PermissionModel
type PermissionMapper struct{}

func NewPermissionMapper() *PermissionMapper {
	return &PermissionMapper{}
}

// ToDomain converts PostgreSQL PermissionModel to domain Permission
func (m *PermissionMapper) ToDomain(model postgres.PermissionModel) permission.Permission {
	return permission.Permission{
		ID:        model.ID,
		Name:      model.Name,
		Resource:  model.Resource,
		Action:    model.Action,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// FromDomain converts domain Permission to PostgreSQL PermissionModel
func (m *PermissionMapper) FromDomain(p permission.Permission) postgres.PermissionModel {
	return postgres.PermissionModel{
		ID:        p.ID,
		Name:      p.Name,
		Resource:  p.Resource,
		Action:    p.Action,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

