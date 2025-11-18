package notification

// DomainEvent represents a domain event that occurred within the notification domain
type DomainEvent interface {
	Type() string
}

// NotificationCreatedEvent is a domain event emitted when a notification is created
type NotificationCreatedEvent struct {
	NotificationID string
	UserID         string
	NotificationType Type
	Title          string
	Body           string
	Metadata       map[string]string
	CreatedAt      string
}

func (e NotificationCreatedEvent) Type() string {
	return "NotificationCreated"
}

// NotificationReadEvent is a domain event emitted when a notification is marked as read
type NotificationReadEvent struct {
	NotificationID string
	UserID         string
	ReadAt         string
}

func (e NotificationReadEvent) Type() string {
	return "NotificationRead"
}

