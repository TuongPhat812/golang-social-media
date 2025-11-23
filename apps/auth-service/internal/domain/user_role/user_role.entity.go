package user_role

import (
	"time"
)

// UserRole represents the relationship between a user and a role
type UserRole struct {
	UserID    string
	RoleID    string
	CreatedAt time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Assign assigns a role to a user and adds a domain event
func (ur *UserRole) Assign() {
	ur.addEvent(UserRoleAssignedEvent{
		UserID:    ur.UserID,
		RoleID:    ur.RoleID,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// Revoke revokes a role from a user and adds a domain event
func (ur *UserRole) Revoke() {
	ur.addEvent(UserRoleRevokedEvent{
		UserID:    ur.UserID,
		RoleID:    ur.RoleID,
		RevokedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// Events returns all domain events
func (ur UserRole) Events() []DomainEvent {
	return ur.events
}

// ClearEvents clears all domain events
func (ur *UserRole) ClearEvents() {
	ur.events = nil
}

// addEvent adds a domain event (internal method)
func (ur *UserRole) addEvent(event DomainEvent) {
	ur.events = append(ur.events, event)
}

