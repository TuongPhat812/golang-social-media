package memory

import (
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/permission"
	pkgerrors "golang-social-media/pkg/errors"
)

var (
	ErrPermissionNotFound      = pkgerrors.NewNotFoundError("permission_not_found")
	ErrPermissionAlreadyExists  = pkgerrors.NewConflictError("permission_already_exists")
)

type PermissionRepository struct {
	mu         sync.RWMutex
	byID       map[string]permission.Permission
	byResourceAction map[string]permission.Permission // key: "resource:action"
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		byID:              make(map[string]permission.Permission),
		byResourceAction:  make(map[string]permission.Permission),
	}
}

func (r *PermissionRepository) key(resource, action string) string {
	return resource + ":" + action
}

func (r *PermissionRepository) Create(perm permission.Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(perm.Resource, perm.Action)
	if _, exists := r.byResourceAction[key]; exists {
		return ErrPermissionAlreadyExists
	}

	r.byID[perm.ID] = perm
	r.byResourceAction[key] = perm
	return nil
}

func (r *PermissionRepository) GetByID(id string) (permission.Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perm, ok := r.byID[id]
	if !ok {
		return permission.Permission{}, ErrPermissionNotFound
	}
	return perm, nil
}

func (r *PermissionRepository) GetByResourceAction(resource, action string) (permission.Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.key(resource, action)
	perm, ok := r.byResourceAction[key]
	if !ok {
		return permission.Permission{}, ErrPermissionNotFound
	}
	return perm, nil
}

func (r *PermissionRepository) List() ([]permission.Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	permissions := make([]permission.Permission, 0, len(r.byID))
	for _, perm := range r.byID {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

