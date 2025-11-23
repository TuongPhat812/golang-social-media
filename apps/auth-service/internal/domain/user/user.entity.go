package user

import (
	"strings"
	"time"

	"golang-social-media/pkg/errors"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	UpdatedAt time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the user
func (u User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.NewValidationError(errors.CodeEmailRequired, nil)
	}
	if !strings.Contains(u.Email, "@") {
		return errors.NewValidationError(errors.CodeEmailInvalid, nil)
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.NewValidationError(errors.CodePasswordRequired, nil)
	}
	if len(u.Password) < 6 {
		return errors.NewValidationError(errors.CodePasswordTooShort, nil)
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.NewValidationError(errors.CodeNameRequired, nil)
	}
	return nil
}

// ValidatePassword validates password rules
func (u User) ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.NewValidationError(errors.CodePasswordRequired, nil)
	}
	if len(password) < 6 {
		return errors.NewValidationError(errors.CodePasswordTooShort, nil)
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

// UpdateProfile updates user profile and adds a domain event
func (u *User) UpdateProfile(name string) {
	if strings.TrimSpace(name) == "" {
		return // Invalid, should be validated before calling
	}
	oldName := u.Name
	u.Name = name
	u.UpdatedAt = time.Now().UTC()

	u.addEvent(UserProfileUpdatedEvent{
		UserID:    u.ID,
		OldName:   oldName,
		NewName:   name,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// ChangePassword changes user password and adds a domain event
func (u *User) ChangePassword(newPassword string) {
	u.Password = newPassword
	u.UpdatedAt = time.Now().UTC()

	u.addEvent(UserPasswordChangedEvent{
		UserID:    u.ID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
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
