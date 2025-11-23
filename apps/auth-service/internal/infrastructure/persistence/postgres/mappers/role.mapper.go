package mappers

import (
	"golang-social-media/apps/auth-service/internal/domain/role"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
)

// RoleMapper maps between domain Role and PostgreSQL RoleModel
type RoleMapper struct{}

func NewRoleMapper() *RoleMapper {
	return &RoleMapper{}
}

// ToDomain converts PostgreSQL RoleModel to domain Role
func (m *RoleMapper) ToDomain(model postgres.RoleModel) role.Role {
	return role.Role{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// FromDomain converts domain Role to PostgreSQL RoleModel
func (m *RoleMapper) FromDomain(r role.Role) postgres.RoleModel {
	return postgres.RoleModel{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

