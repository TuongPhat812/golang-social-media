package query

import (
	"context"
	"errors"
	"testing"

	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/domain/permission"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPermissionRepository is a mock implementation for testing
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(p permission.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByID(id string) (permission.Permission, error) {
	args := m.Called(id)
	return args.Get(0).(permission.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByResourceAction(resource, action string) (permission.Permission, error) {
	args := m.Called(resource, action)
	return args.Get(0).(permission.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetAll() ([]permission.Permission, error) {
	args := m.Called()
	return args.Get(0).([]permission.Permission), args.Error(1)
}

// MockUserRoleRepository is a mock implementation for testing
type MockCheckPermissionUserRoleRepository struct {
	mock.Mock
}

func (m *MockCheckPermissionUserRoleRepository) Create(ur interface{}) error {
	args := m.Called(ur)
	return args.Error(0)
}

func (m *MockCheckPermissionUserRoleRepository) Delete(userID, roleID string) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockCheckPermissionUserRoleRepository) GetUserRoles(userID string) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

// MockRolePermissionRepository is a mock implementation for testing
type MockRolePermissionRepository struct {
	mock.Mock
}

func (m *MockRolePermissionRepository) Create(rp interface{}) error {
	args := m.Called(rp)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) Delete(roleID, permissionID string) error {
	args := m.Called(roleID, permissionID)
	return args.Error(0)
}

func (m *MockRolePermissionRepository) HasPermission(roleID, permissionID string) (bool, error) {
	args := m.Called(roleID, permissionID)
	return args.Bool(0), args.Error(1)
}

func TestCheckPermissionQuery_Execute(t *testing.T) {
	ctx := context.Background()
	req := contracts.CheckPermissionQueryRequest{
		UserID:   "user-1",
		Resource: "chat",
		Action:   "create",
	}

	t.Run("User Has Permission", func(t *testing.T) {
		mockPermRepo := new(MockPermissionRepository)
		mockUserRoleRepo := new(MockCheckPermissionUserRoleRepository)
		mockRolePermRepo := new(MockRolePermissionRepository)

		testPermission := permission.Permission{
			ID:       "perm-1",
			Name:     "Create Chat",
			Resource: "chat",
			Action:   "create",
		}
		roleIDs := []string{"role-1", "role-2"}

		mockPermRepo.On("GetByResourceAction", "chat", "create").Return(testPermission, nil)
		mockUserRoleRepo.On("GetUserRoles", "user-1").Return(roleIDs, nil)
		mockRolePermRepo.On("HasPermission", "role-1", "perm-1").Return(true, nil)

		query := NewCheckPermissionQuery(mockUserRoleRepo, mockRolePermRepo, mockPermRepo)
		resp, err := query.Execute(ctx, req)

		assert.Nil(t, err)
		assert.True(t, resp.HasPermission)

		mockPermRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
		mockRolePermRepo.AssertExpectations(t)
	})

	t.Run("User Does Not Have Permission", func(t *testing.T) {
		mockPermRepo := new(MockPermissionRepository)
		mockUserRoleRepo := new(MockCheckPermissionUserRoleRepository)
		mockRolePermRepo := new(MockRolePermissionRepository)

		testPermission := permission.Permission{
			ID:       "perm-1",
			Name:     "Create Chat",
			Resource: "chat",
			Action:   "create",
		}
		roleIDs := []string{"role-1"}

		mockPermRepo.On("GetByResourceAction", "chat", "create").Return(testPermission, nil)
		mockUserRoleRepo.On("GetUserRoles", "user-1").Return(roleIDs, nil)
		mockRolePermRepo.On("HasPermission", "role-1", "perm-1").Return(false, nil)

		query := NewCheckPermissionQuery(mockUserRoleRepo, mockRolePermRepo, mockPermRepo)
		resp, err := query.Execute(ctx, req)

		assert.Nil(t, err)
		assert.False(t, resp.HasPermission)

		mockPermRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
		mockRolePermRepo.AssertExpectations(t)
	})

	t.Run("Permission Not Found", func(t *testing.T) {
		mockPermRepo := new(MockPermissionRepository)
		mockUserRoleRepo := new(MockCheckPermissionUserRoleRepository)
		mockRolePermRepo := new(MockRolePermissionRepository)

		permErr := errors.New("permission not found")
		mockPermRepo.On("GetByResourceAction", "chat", "create").Return(permission.Permission{}, permErr)

		query := NewCheckPermissionQuery(mockUserRoleRepo, mockRolePermRepo, mockPermRepo)
		resp, err := query.Execute(ctx, req)

		assert.Nil(t, err)
		assert.False(t, resp.HasPermission)

		mockPermRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertNotCalled(t, "GetUserRoles", mock.Anything)
		mockRolePermRepo.AssertNotCalled(t, "HasPermission", mock.Anything, mock.Anything)
	})

	t.Run("GetUserRoles Fails", func(t *testing.T) {
		mockPermRepo := new(MockPermissionRepository)
		mockUserRoleRepo := new(MockCheckPermissionUserRoleRepository)
		mockRolePermRepo := new(MockRolePermissionRepository)

		testPermission := permission.Permission{
			ID:       "perm-1",
			Name:     "Create Chat",
			Resource: "chat",
			Action:   "create",
		}
		userRoleErr := errors.New("database error")

		mockPermRepo.On("GetByResourceAction", "chat", "create").Return(testPermission, nil)
		mockUserRoleRepo.On("GetUserRoles", "user-1").Return([]string{}, userRoleErr)

		query := NewCheckPermissionQuery(mockUserRoleRepo, mockRolePermRepo, mockPermRepo)
		resp, err := query.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.False(t, resp.HasPermission)

		mockPermRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
		mockRolePermRepo.AssertNotCalled(t, "HasPermission", mock.Anything, mock.Anything)
	})

	t.Run("HasPermission Check Fails (Continues)", func(t *testing.T) {
		mockPermRepo := new(MockPermissionRepository)
		mockUserRoleRepo := new(MockCheckPermissionUserRoleRepository)
		mockRolePermRepo := new(MockRolePermissionRepository)

		testPermission := permission.Permission{
			ID:       "perm-1",
			Name:     "Create Chat",
			Resource: "chat",
			Action:   "create",
		}
		roleIDs := []string{"role-1", "role-2"}

		mockPermRepo.On("GetByResourceAction", "chat", "create").Return(testPermission, nil)
		mockUserRoleRepo.On("GetUserRoles", "user-1").Return(roleIDs, nil)
		mockRolePermRepo.On("HasPermission", "role-1", "perm-1").Return(false, errors.New("database error"))
		mockRolePermRepo.On("HasPermission", "role-2", "perm-1").Return(true, nil)

		query := NewCheckPermissionQuery(mockUserRoleRepo, mockRolePermRepo, mockPermRepo)
		resp, err := query.Execute(ctx, req)

		assert.Nil(t, err)
		assert.True(t, resp.HasPermission) // Should find permission via role-2

		mockPermRepo.AssertExpectations(t)
		mockUserRoleRepo.AssertExpectations(t)
		mockRolePermRepo.AssertExpectations(t)
	})
}

