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

// UserProfileUpdatedEvent is a domain event emitted when user profile is updated
type UserProfileUpdatedEvent struct {
	UserID    string
	OldName   string
	NewName   string
	UpdatedAt string
}

func (e UserProfileUpdatedEvent) Type() string {
	return "UserProfileUpdated"
}

// UserPasswordChangedEvent is a domain event emitted when user password is changed
type UserPasswordChangedEvent struct {
	UserID    string
	UpdatedAt string
}

func (e UserPasswordChangedEvent) Type() string {
	return "UserPasswordChanged"
}
