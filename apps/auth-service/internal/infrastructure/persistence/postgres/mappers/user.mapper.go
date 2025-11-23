package mappers

import (
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
)

// UserMapper maps between domain User and PostgreSQL UserModel
type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToDomain converts PostgreSQL UserModel to domain User
func (m *UserMapper) ToDomain(model postgres.UserModel) user.User {
	return user.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Name:      model.Name,
		UpdatedAt: model.UpdatedAt,
	}
}

// FromDomain converts domain User to PostgreSQL UserModel
func (m *UserMapper) FromDomain(u user.User) postgres.UserModel {
	return postgres.UserModel{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		UpdatedAt: u.UpdatedAt,
	}
}

