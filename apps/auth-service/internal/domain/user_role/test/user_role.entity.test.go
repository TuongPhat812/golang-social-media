package user_role

import (
	"testing"
	"time"
)

func TestUserRole_Assign(t *testing.T) {
	userRole := &UserRole{
		UserID:    "user-1",
		RoleID:    "role-1",
		CreatedAt: time.Now(),
	}

	userRole.Assign()

	events := userRole.Events()
	if len(events) != 1 {
		t.Fatalf("UserRole.Assign() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(UserRoleAssignedEvent)
	if !ok {
		t.Fatalf("UserRole.Assign() event type = %T, want UserRoleAssignedEvent", events[0])
	}

	if event.UserID != userRole.UserID {
		t.Errorf("UserRoleAssignedEvent.UserID = %v, want %v", event.UserID, userRole.UserID)
	}
	if event.RoleID != userRole.RoleID {
		t.Errorf("UserRoleAssignedEvent.RoleID = %v, want %v", event.RoleID, userRole.RoleID)
	}
	if event.CreatedAt == "" {
		t.Error("UserRoleAssignedEvent.CreatedAt should not be empty")
	}
}

func TestUserRole_Revoke(t *testing.T) {
	userRole := &UserRole{
		UserID:    "user-1",
		RoleID:    "role-1",
		CreatedAt: time.Now(),
	}

	userRole.Revoke()

	events := userRole.Events()
	if len(events) != 1 {
		t.Fatalf("UserRole.Revoke() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(UserRoleRevokedEvent)
	if !ok {
		t.Fatalf("UserRole.Revoke() event type = %T, want UserRoleRevokedEvent", events[0])
	}

	if event.UserID != userRole.UserID {
		t.Errorf("UserRoleRevokedEvent.UserID = %v, want %v", event.UserID, userRole.UserID)
	}
	if event.RoleID != userRole.RoleID {
		t.Errorf("UserRoleRevokedEvent.RoleID = %v, want %v", event.RoleID, userRole.RoleID)
	}
	if event.RevokedAt == "" {
		t.Error("UserRoleRevokedEvent.RevokedAt should not be empty")
	}
}

func TestUserRole_ClearEvents(t *testing.T) {
	userRole := &UserRole{
		UserID: "user-1",
		RoleID: "role-1",
	}

	userRole.Assign()
	if len(userRole.Events()) != 1 {
		t.Fatal("Expected 1 event after Assign()")
	}

	userRole.ClearEvents()
	if len(userRole.Events()) != 0 {
		t.Errorf("UserRole.ClearEvents() should clear all events, got %d", len(userRole.Events()))
	}
}

