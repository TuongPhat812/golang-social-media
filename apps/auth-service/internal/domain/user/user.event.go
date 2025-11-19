package user

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// UserCreatedEvent is a domain event emitted when a user is created
type UserCreatedEvent struct {
	UserID    string
	Email     string
	Name      string
	CreatedAt string
}

func (e UserCreatedEvent) Type() string {
	return "UserCreated"
}

