package mappers

import (
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
)

// UserMapper maps between domain User and persistence models
type UserMapper struct{}

// NewUserMapper creates a new UserMapper
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToDomain converts a memory user representation to domain User
// Note: For memory repository, we store domain entities directly,
// but this mapper provides a consistent interface for future database implementations
func (m *UserMapper) ToDomain(u user.User) user.User {
	// For memory repository, this is a pass-through
	// In a real database implementation, this would convert from DB model
	return u
}

// FromDomain converts a domain User to memory representation
func (m *UserMapper) FromDomain(u user.User) user.User {
	// For memory repository, this is a pass-through
	return u
}

// ToDomainList converts a slice of users to domain Users
func (m *UserMapper) ToDomainList(users []user.User) []user.User {
	result := make([]user.User, len(users))
	for i, u := range users {
		result[i] = m.ToDomain(u)
	}
	return result
}

// UserModelMapper is a placeholder for future database implementations
// When migrating to PostgreSQL, this will map between UserModel and domain User
type UserModelMapper struct{}

func NewUserModelMapper() *UserModelMapper {
	return &UserModelMapper{}
}

// ToDomain converts a database model to domain User
// This will be implemented when we add PostgreSQL support
func (m *UserModelMapper) ToDomain(model interface{}) user.User {
	// TODO: Implement when adding database support
	panic("not implemented")
}

// FromDomain converts a domain User to database model
func (m *UserModelMapper) FromDomain(u user.User) interface{} {
	// TODO: Implement when adding database support
	panic("not implemented")
}

