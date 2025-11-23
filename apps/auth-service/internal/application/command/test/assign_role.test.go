package command

import (
	"context"
	"errors"
	"testing"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/role"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
	pkgerrors "golang-social-media/pkg/errors"

	"github.com/stretchr/testify/mock"
)

// MockUserRoleRepository is a mock implementation for testing
type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ur interface{}) error {
	args := m.Called(ur)
	return args.Error(0)
}

func (m *MockUserRoleRepository) Delete(userID, roleID string) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) GetUserRoles(userID string) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

// MockRoleRepository is a mock implementation for testing
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(r role.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(id string) (role.Role, error) {
	args := m.Called(id)
	return args.Get(0).(role.Role), args.Error(1)
}

func (m *MockRoleRepository) GetAll() ([]role.Role, error) {
	args := m.Called()
	return args.Get(0).([]role.Role), args.Error(1)
}

func TestAssignRoleCommand_Execute(t *testing.T) {
	ctx := context.Background()
	req := contracts.AssignRoleCommandRequest{
		UserID: "user-1",
		RoleID: "role-1",
	}

	t.Run("Successful Role Assignment", func(t *testing.T) {
		mockUserRoleRepo := new(MockUserRoleRepository)
		mockRoleRepo := new(MockRoleRepository)

		testRole := role.Role{
			ID:   "role-1",
			Name: "Admin",
		}

		mockRoleRepo.On("GetByID", "role-1").Return(testRole, nil)
		mockUserRoleRepo.On("Create", mock.AnythingOfType("user_role.UserRole")).Return(nil)

		cmd := NewAssignRoleCommand(mockUserRoleRepo, mockRoleRepo)
		err := cmd.Execute(ctx, req)

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		mockRoleRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
	})

	t.Run("Role Not Found", func(t *testing.T) {
		mockUserRoleRepo := new(MockUserRoleRepository)
		mockRoleRepo := new(MockRoleRepository)

		roleErr := pkgerrors.NewNotFoundError("role_not_found")
		mockRoleRepo.On("GetByID", "role-1").Return(role.Role{}, roleErr)

		cmd := NewAssignRoleCommand(mockUserRoleRepo, mockRoleRepo)
		err := cmd.Execute(ctx, req)

		if err == nil {
			t.Error("Execute() should return error when role not found")
		}

		mockRoleRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Repository Create Fails", func(t *testing.T) {
		mockUserRoleRepo := new(MockUserRoleRepository)
		mockRoleRepo := new(MockRoleRepository)

		testRole := role.Role{
			ID:   "role-1",
			Name: "Admin",
		}

		repoErr := errors.New("database error")
		mockRoleRepo.On("GetByID", "role-1").Return(testRole, nil)
		mockUserRoleRepo.On("Create", mock.AnythingOfType("user_role.UserRole")).Return(repoErr)

		cmd := NewAssignRoleCommand(mockUserRoleRepo, mockRoleRepo)
		err := cmd.Execute(ctx, req)

		if err == nil {
			t.Error("Execute() should return error when repository create fails")
		}

		mockRoleRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
	})
}

