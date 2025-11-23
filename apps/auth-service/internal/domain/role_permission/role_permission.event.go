package role_permission

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// RolePermissionAssignedEvent is a domain event emitted when a permission is assigned to a role
type RolePermissionAssignedEvent struct {
	RoleID       string
	PermissionID string
	CreatedAt    string
}

func (e RolePermissionAssignedEvent) Type() string {
	return "RolePermissionAssigned"
}

// RolePermissionRevokedEvent is a domain event emitted when a permission is revoked from a role
type RolePermissionRevokedEvent struct {
	RoleID       string
	PermissionID string
	RevokedAt    string
}

func (e RolePermissionRevokedEvent) Type() string {
	return "RolePermissionRevoked"
}

