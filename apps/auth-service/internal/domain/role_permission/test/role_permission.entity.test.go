package role_permission

import (
	"testing"
	"time"
)

func TestRolePermission_Assign(t *testing.T) {
	rolePerm := &RolePermission{
		RoleID:       "role-1",
		PermissionID: "perm-1",
		CreatedAt:    time.Now(),
	}

	rolePerm.Assign()

	events := rolePerm.Events()
	if len(events) != 1 {
		t.Fatalf("RolePermission.Assign() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(RolePermissionAssignedEvent)
	if !ok {
		t.Fatalf("RolePermission.Assign() event type = %T, want RolePermissionAssignedEvent", events[0])
	}

	if event.RoleID != rolePerm.RoleID {
		t.Errorf("RolePermissionAssignedEvent.RoleID = %v, want %v", event.RoleID, rolePerm.RoleID)
	}
	if event.PermissionID != rolePerm.PermissionID {
		t.Errorf("RolePermissionAssignedEvent.PermissionID = %v, want %v", event.PermissionID, rolePerm.PermissionID)
	}
	if event.CreatedAt == "" {
		t.Error("RolePermissionAssignedEvent.CreatedAt should not be empty")
	}
}

func TestRolePermission_Revoke(t *testing.T) {
	rolePerm := &RolePermission{
		RoleID:       "role-1",
		PermissionID: "perm-1",
		CreatedAt:    time.Now(),
	}

	rolePerm.Revoke()

	events := rolePerm.Events()
	if len(events) != 1 {
		t.Fatalf("RolePermission.Revoke() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(RolePermissionRevokedEvent)
	if !ok {
		t.Fatalf("RolePermission.Revoke() event type = %T, want RolePermissionRevokedEvent", events[0])
	}

	if event.RoleID != rolePerm.RoleID {
		t.Errorf("RolePermissionRevokedEvent.RoleID = %v, want %v", event.RoleID, rolePerm.RoleID)
	}
	if event.PermissionID != rolePerm.PermissionID {
		t.Errorf("RolePermissionRevokedEvent.PermissionID = %v, want %v", event.PermissionID, rolePerm.PermissionID)
	}
	if event.RevokedAt == "" {
		t.Error("RolePermissionRevokedEvent.RevokedAt should not be empty")
	}
}

func TestRolePermission_ClearEvents(t *testing.T) {
	rolePerm := &RolePermission{
		RoleID:       "role-1",
		PermissionID: "perm-1",
	}

	rolePerm.Assign()
	if len(rolePerm.Events()) != 1 {
		t.Fatal("Expected 1 event after Assign()")
	}

	rolePerm.ClearEvents()
	if len(rolePerm.Events()) != 0 {
		t.Errorf("RolePermission.ClearEvents() should clear all events, got %d", len(rolePerm.Events()))
	}
}

