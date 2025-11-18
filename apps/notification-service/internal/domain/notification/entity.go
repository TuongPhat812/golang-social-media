package notification

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type Type string

const (
	TypeWelcome     Type = "welcome"
	TypeChatMessage Type = "chat_message"
)

type Notification struct {
	ID        gocql.UUID
	UserID    string
	Type      Type
	Title     string
	Body      string
	Metadata  map[string]string
	CreatedAt time.Time
	ReadAt    *time.Time // nil if unread

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate performs business logic validation on the notification
// This method should be called before persisting to database or publishing events
func (n Notification) Validate() error {
	if n.UserID == "" {
		return errors.New("user_id is required")
	}

	if n.Type == "" {
		return errors.New("type is required")
	}

	if n.Type != TypeWelcome && n.Type != TypeChatMessage {
		return errors.New("invalid notification type")
	}

	if n.Title == "" {
		return errors.New("title is required")
	}

	if n.Body == "" {
		return errors.New("body is required")
	}

	if n.CreatedAt.IsZero() {
		return errors.New("created_at is required")
	}

	// Add more business rules here as needed
	// Example: if n.Type == TypeAdultContent && userAge < 18 {
	//     return errors.New("adult content requires user to be 18+")
	// }

	return nil
}

// Create is a domain method that creates a notification and adds a domain event
// This method should be called to create a notification with proper domain logic
func (n *Notification) Create() {
	// Add domain event when notification is created
	n.addEvent(NotificationCreatedEvent{
		NotificationID: n.ID.String(),
		UserID:         n.UserID,
		NotificationType: n.Type,
		Title:          n.Title,
		Body:           n.Body,
		Metadata:       n.Metadata,
		CreatedAt:      n.CreatedAt.Format(time.RFC3339),
	})
}

// addEvent adds a domain event to the notification
func (n *Notification) addEvent(event DomainEvent) {
	if n.events == nil {
		n.events = make([]DomainEvent, 0)
	}
	n.events = append(n.events, event)
}

// Events returns all domain events that occurred on this notification
func (n *Notification) Events() []DomainEvent {
	if n.events == nil {
		return []DomainEvent{}
	}
	return n.events
}

// ClearEvents clears all domain events from the notification
// This should be called after events have been dispatched
func (n *Notification) ClearEvents() {
	n.events = nil
}

// MarkAsRead marks the notification as read and adds a domain event
// This method should be called to mark a notification as read with proper domain logic
func (n *Notification) MarkAsRead() error {
	if n.ReadAt != nil {
		return errors.New("notification is already read")
	}

	now := time.Now().UTC()
	n.ReadAt = &now

	// Add domain event when notification is marked as read
	n.addEvent(NotificationReadEvent{
		NotificationID: n.ID.String(),
		UserID:         n.UserID,
		ReadAt:         now.Format(time.RFC3339),
	})

	return nil
}

// IsRead returns true if the notification has been read
func (n Notification) IsRead() bool {
	return n.ReadAt != nil
}
