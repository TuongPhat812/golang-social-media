package role

import (
	"strings"
	"time"

	"golang-social-media/pkg/errors"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the role
func (r Role) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.NewValidationError(errors.CodeNameRequired, nil)
	}
	if len(r.Name) < 2 {
		return errors.NewValidationError(errors.CodeNameTooShort, nil)
	}
	if len(r.Name) > 50 {
		return errors.NewValidationError(errors.CodeNameTooLong, nil)
	}
	return nil
}

// Create is a domain method that creates a role and adds a domain event
func (r *Role) Create() {
	r.addEvent(RoleCreatedEvent{
		RoleID:    r.ID,
		Name:      r.Name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// Update updates role information
func (r *Role) Update(name, description string) {
	oldName := r.Name
	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now().UTC()

	r.addEvent(RoleUpdatedEvent{
		RoleID:    r.ID,
		OldName:   oldName,
		NewName:   name,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// Events returns all domain events
func (r Role) Events() []DomainEvent {
	return r.events
}

// ClearEvents clears all domain events
func (r *Role) ClearEvents() {
	r.events = nil
}

// addEvent adds a domain event (internal method)
func (r *Role) addEvent(event DomainEvent) {
	r.events = append(r.events, event)
}

