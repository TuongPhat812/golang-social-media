package permission

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// PermissionCreatedEvent is a domain event emitted when a permission is created
type PermissionCreatedEvent struct {
	PermissionID string
	Name         string
	Resource     string
	Action       string
	CreatedAt    string
}

func (e PermissionCreatedEvent) Type() string {
	return "PermissionCreated"
}

