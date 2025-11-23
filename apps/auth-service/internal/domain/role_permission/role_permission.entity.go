package role_permission

import (
	"time"
)

// RolePermission represents the relationship between a role and a permission
type RolePermission struct {
	RoleID       string
	PermissionID string
	CreatedAt    time.Time

	// Domain events (internal, not persisted)
	events []DomainEvent
}

// Assign assigns a permission to a role and adds a domain event
func (rp *RolePermission) Assign() {
	rp.addEvent(RolePermissionAssignedEvent{
		RoleID:       rp.RoleID,
		PermissionID: rp.PermissionID,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	})
}

// Revoke revokes a permission from a role and adds a domain event
func (rp *RolePermission) Revoke() {
	rp.addEvent(RolePermissionRevokedEvent{
		RoleID:       rp.RoleID,
		PermissionID: rp.PermissionID,
		RevokedAt:    time.Now().UTC().Format(time.RFC3339),
	})
}

// Events returns all domain events
func (rp RolePermission) Events() []DomainEvent {
	return rp.events
}

// ClearEvents clears all domain events
func (rp *RolePermission) ClearEvents() {
	rp.events = nil
}

// addEvent adds a domain event (internal method)
func (rp *RolePermission) addEvent(event DomainEvent) {
	rp.events = append(rp.events, event)
}

