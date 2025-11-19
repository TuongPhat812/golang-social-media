package user

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID       string
	Email    string
	Password string
	Name     string

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the user
func (u User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("email must be valid")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password cannot be empty")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

// Create is a domain method that creates a user and adds a domain event
func (u *User) Create() {
	u.addEvent(UserCreatedEvent{
		UserID:    u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// Events returns all domain events
func (u User) Events() []DomainEvent {
	return u.events
}

// ClearEvents clears all domain events
func (u *User) ClearEvents() {
	u.events = nil
}

// addEvent adds a domain event (internal method)
func (u *User) addEvent(event DomainEvent) {
	u.events = append(u.events, event)
}
