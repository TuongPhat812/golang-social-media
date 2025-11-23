package memory

import (
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/user_role"
)

type UserRoleRepository struct {
	mu        sync.RWMutex
	byUserID  map[string][]string // userID -> []roleID
	byRoleID  map[string][]string // roleID -> []userID
	userRoles map[string]user_role.UserRole // key: "userID:roleID"
}

func NewUserRoleRepository() *UserRoleRepository {
	return &UserRoleRepository{
		byUserID:  make(map[string][]string),
		byRoleID:  make(map[string][]string),
		userRoles: make(map[string]user_role.UserRole),
	}
}

func (r *UserRoleRepository) key(userID, roleID string) string {
	return userID + ":" + roleID
}

func (r *UserRoleRepository) Create(userRole user_role.UserRole) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(userRole.UserID, userRole.RoleID)
	if _, exists := r.userRoles[key]; exists {
		// Already assigned, skip
		return nil
	}

	r.userRoles[key] = userRole
	r.byUserID[userRole.UserID] = append(r.byUserID[userRole.UserID], userRole.RoleID)
	r.byRoleID[userRole.RoleID] = append(r.byRoleID[userRole.RoleID], userRole.UserID)
	return nil
}

func (r *UserRoleRepository) GetUserRoles(userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roleIDs, ok := r.byUserID[userID]
	if !ok {
		return []string{}, nil
	}
	return roleIDs, nil
}

func (r *UserRoleRepository) GetRoleUsers(roleID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userIDs, ok := r.byRoleID[roleID]
	if !ok {
		return []string{}, nil
	}
	return userIDs, nil
}

func (r *UserRoleRepository) Delete(userID, roleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(userID, roleID)
	if _, exists := r.userRoles[key]; !exists {
		// Not assigned, skip
		return nil
	}

	delete(r.userRoles, key)

	// Remove from byUserID
	roleIDs := r.byUserID[userID]
	newRoleIDs := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		if id != roleID {
			newRoleIDs = append(newRoleIDs, id)
		}
	}
	r.byUserID[userID] = newRoleIDs

	// Remove from byRoleID
	userIDs := r.byRoleID[roleID]
	newUserIDs := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		if id != userID {
			newUserIDs = append(newUserIDs, id)
		}
	}
	r.byRoleID[roleID] = newUserIDs

	return nil
}

func (r *UserRoleRepository) HasRole(userID, roleID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.key(userID, roleID)
	_, exists := r.userRoles[key]
	return exists, nil
}

