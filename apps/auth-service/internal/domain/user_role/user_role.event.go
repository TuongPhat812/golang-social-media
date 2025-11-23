package user_role

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// UserRoleAssignedEvent is a domain event emitted when a role is assigned to a user
type UserRoleAssignedEvent struct {
	UserID    string
	RoleID    string
	CreatedAt string
}

func (e UserRoleAssignedEvent) Type() string {
	return "UserRoleAssigned"
}

// UserRoleRevokedEvent is a domain event emitted when a role is revoked from a user
type UserRoleRevokedEvent struct {
	UserID    string
	RoleID    string
	RevokedAt string
}

func (e UserRoleRevokedEvent) Type() string {
	return "UserRoleRevoked"
}

