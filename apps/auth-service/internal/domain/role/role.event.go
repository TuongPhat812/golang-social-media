package role

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// RoleCreatedEvent is a domain event emitted when a role is created
type RoleCreatedEvent struct {
	RoleID    string
	Name      string
	CreatedAt string
}

func (e RoleCreatedEvent) Type() string {
	return "RoleCreated"
}

// RoleUpdatedEvent is a domain event emitted when a role is updated
type RoleUpdatedEvent struct {
	RoleID    string
	OldName   string
	NewName   string
	UpdatedAt string
}

func (e RoleUpdatedEvent) Type() string {
	return "RoleUpdated"
}

