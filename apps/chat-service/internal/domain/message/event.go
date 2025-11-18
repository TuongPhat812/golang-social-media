package message

// DomainEvent represents a domain event interface
type DomainEvent interface {
	Type() string
}

// MessageCreatedEvent is a domain event emitted when a message is created
type MessageCreatedEvent struct {
	MessageID  string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  string
}

func (e MessageCreatedEvent) Type() string {
	return "MessageCreated"
}
