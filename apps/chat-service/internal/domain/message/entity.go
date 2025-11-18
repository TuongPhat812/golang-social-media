package message

import (
	"errors"
	"strings"
	"time"
)

type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the message
func (m Message) Validate() error {
	if strings.TrimSpace(m.SenderID) == "" {
		return errors.New("sender ID cannot be empty")
	}
	if strings.TrimSpace(m.ReceiverID) == "" {
		return errors.New("receiver ID cannot be empty")
	}
	if m.SenderID == m.ReceiverID {
		return errors.New("sender and receiver cannot be the same")
	}
	if strings.TrimSpace(m.Content) == "" {
		return errors.New("message content cannot be empty")
	}
	if len(m.Content) > 5000 {
		return errors.New("message content cannot exceed 5000 characters")
	}
	return nil
}

// Create is a domain method that creates a message and adds a domain event
func (m *Message) Create() {
	m.addEvent(MessageCreatedEvent{
		MessageID:  m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	})
}

// Events returns all domain events
func (m Message) Events() []DomainEvent {
	return m.events
}

// ClearEvents clears all domain events
func (m *Message) ClearEvents() {
	m.events = nil
}

// addEvent adds a domain event (internal method)
func (m *Message) addEvent(event DomainEvent) {
	m.events = append(m.events, event)
}
