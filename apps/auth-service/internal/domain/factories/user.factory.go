package factories

import (
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/pkg/random"
)

// UserFactoryImpl implements UserFactory interface
type UserFactoryImpl struct {
	idGenerator func() string
}

var _ UserFactory = (*UserFactoryImpl)(nil)

// NewUserFactory creates a new UserFactoryImpl
func NewUserFactory() UserFactory {
	return &UserFactoryImpl{
		idGenerator: func() string {
			return "user-" + random.String(8)
		},
	}
}

// NewUserFactoryWithIDGenerator creates a UserFactoryImpl with custom ID generator
func NewUserFactoryWithIDGenerator(idGenerator func() string) UserFactory {
	return &UserFactoryImpl{
		idGenerator: idGenerator,
	}
}

// CreateUser creates a new User with proper initialization
// This factory encapsulates the complex creation logic
func (f *UserFactoryImpl) CreateUser(email, password, name string) (*user.User, error) {
	if email == "" {
		return nil, &UserFactoryError{Message: "email cannot be empty"}
	}
	if password == "" {
		return nil, &UserFactoryError{Message: "password cannot be empty"}
	}
	if name == "" {
		return nil, &UserFactoryError{Message: "name cannot be empty"}
	}

	userModel := &user.User{
		ID:       f.idGenerator(),
		Email:    email,
		Password: password,
		Name:     name,
	}

	// Validate the created user
	if err := userModel.Validate(); err != nil {
		return nil, &UserFactoryError{
			Message: "failed to validate user",
			Cause:   err,
		}
	}

	// Domain logic: create user (this adds domain events internally)
	userModel.Create()

	return userModel, nil
}

// UserFactoryError represents an error in user factory
type UserFactoryError struct {
	Message string
	Cause   error
}

func (e *UserFactoryError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *UserFactoryError) Unwrap() error {
	return e.Cause
}

