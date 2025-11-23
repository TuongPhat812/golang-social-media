package memory

import (
	"sync"

	"golang-social-media/apps/auth-service/internal/domain/role"
	pkgerrors "golang-social-media/pkg/errors"
)

var (
	ErrRoleNotFound      = pkgerrors.NewNotFoundError("role_not_found")
	ErrRoleAlreadyExists = pkgerrors.NewConflictError("role_already_exists")
)

type RoleRepository struct {
	mu    sync.RWMutex
	byID  map[string]role.Role
	byName map[string]role.Role
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		byID:   make(map[string]role.Role),
		byName: make(map[string]role.Role),
	}
}

func (r *RoleRepository) Create(roleEntity role.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byName[roleEntity.Name]; exists {
		return ErrRoleAlreadyExists
	}

	r.byID[roleEntity.ID] = roleEntity
	r.byName[roleEntity.Name] = roleEntity
	return nil
}

func (r *RoleRepository) GetByID(id string) (role.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roleEntity, ok := r.byID[id]
	if !ok {
		return role.Role{}, ErrRoleNotFound
	}
	return roleEntity, nil
}

func (r *RoleRepository) GetByName(name string) (role.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roleEntity, ok := r.byName[name]
	if !ok {
		return role.Role{}, ErrRoleNotFound
	}
	return roleEntity, nil
}

func (r *RoleRepository) Update(roleEntity role.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldRole, exists := r.byID[roleEntity.ID]
	if !exists {
		return ErrRoleNotFound
	}

	// If name changed, check if new name already exists
	if oldRole.Name != roleEntity.Name {
		if _, exists := r.byName[roleEntity.Name]; exists {
			return ErrRoleAlreadyExists
		}
		delete(r.byName, oldRole.Name)
	}

	r.byID[roleEntity.ID] = roleEntity
	r.byName[roleEntity.Name] = roleEntity
	return nil
}

func (r *RoleRepository) List() ([]role.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roles := make([]role.Role, 0, len(r.byID))
	for _, roleEntity := range r.byID {
		roles = append(roles, roleEntity)
	}
	return roles, nil
}

