package permission

import (
	"strings"
	"time"

	"golang-social-media/pkg/errors"
)

// Permission represents a permission in the RBAC system
type Permission struct {
	ID       string
	Name     string
	Resource string // e.g., "chat", "user", "notification"
	Action   string // e.g., "create", "read", "update", "delete"
	CreatedAt time.Time
	UpdatedAt time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Validate validates business rules for the permission
func (p Permission) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.NewValidationError(errors.CodeNameRequired, nil)
	}
	if strings.TrimSpace(p.Resource) == "" {
		return errors.NewValidationError("resource_required", nil)
	}
	if strings.TrimSpace(p.Action) == "" {
		return errors.NewValidationError("action_required", nil)
	}
	return nil
}

// Create is a domain method that creates a permission and adds a domain event
func (p *Permission) Create() {
	p.addEvent(PermissionCreatedEvent{
		PermissionID: p.ID,
		Name:         p.Name,
		Resource:     p.Resource,
		Action:       p.Action,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
}

// Events returns all domain events
func (p Permission) Events() []DomainEvent {
	return p.events
}

// ClearEvents clears all domain events
func (p *Permission) ClearEvents() {
	p.events = nil
}

// addEvent adds a domain event (internal method)
func (p *Permission) addEvent(event DomainEvent) {
	p.events = append(p.events, event)
}

