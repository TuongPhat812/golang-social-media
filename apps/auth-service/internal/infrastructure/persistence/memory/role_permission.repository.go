package memory

import (
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/role_permission"
)

type RolePermissionRepository struct {
	mu              sync.RWMutex
	byRoleID        map[string][]string // roleID -> []permissionID
	byPermissionID  map[string][]string // permissionID -> []roleID
	rolePermissions map[string]role_permission.RolePermission // key: "roleID:permissionID"
}

func NewRolePermissionRepository() *RolePermissionRepository {
	return &RolePermissionRepository{
		byRoleID:        make(map[string][]string),
		byPermissionID:  make(map[string][]string),
		rolePermissions: make(map[string]role_permission.RolePermission),
	}
}

func (r *RolePermissionRepository) key(roleID, permissionID string) string {
	return roleID + ":" + permissionID
}

func (r *RolePermissionRepository) Create(rp role_permission.RolePermission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(rp.RoleID, rp.PermissionID)
	if _, exists := r.rolePermissions[key]; exists {
		// Already assigned, skip
		return nil
	}

	r.rolePermissions[key] = rp
	r.byRoleID[rp.RoleID] = append(r.byRoleID[rp.RoleID], rp.PermissionID)
	r.byPermissionID[rp.PermissionID] = append(r.byPermissionID[rp.PermissionID], rp.RoleID)
	return nil
}

func (r *RolePermissionRepository) GetRolePermissions(roleID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permissionIDs, ok := r.byRoleID[roleID]
	if !ok {
		return []string{}, nil
	}
	return permissionIDs, nil
}

func (r *RolePermissionRepository) GetPermissionRoles(permissionID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roleIDs, ok := r.byPermissionID[permissionID]
	if !ok {
		return []string{}, nil
	}
	return roleIDs, nil
}

func (r *RolePermissionRepository) Delete(roleID, permissionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(roleID, permissionID)
	if _, exists := r.rolePermissions[key]; !exists {
		// Not assigned, skip
		return nil
	}

	delete(r.rolePermissions, key)

	// Remove from byRoleID
	permissionIDs := r.byRoleID[roleID]
	newPermissionIDs := make([]string, 0, len(permissionIDs))
	for _, id := range permissionIDs {
		if id != permissionID {
			newPermissionIDs = append(newPermissionIDs, id)
		}
	}
	r.byRoleID[roleID] = newPermissionIDs

	// Remove from byPermissionID
	roleIDs := r.byPermissionID[permissionID]
	newRoleIDs := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		if id != roleID {
			newRoleIDs = append(newRoleIDs, id)
		}
	}
	r.byPermissionID[permissionID] = newRoleIDs

	return nil
}

func (r *RolePermissionRepository) HasPermission(roleID, permissionID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.key(roleID, permissionID)
	_, exists := r.rolePermissions[key]
	return exists, nil
}

