package role

import (
	"testing"
	"time"

	"golang-social-media/pkg/errors"
)

func TestRole_Validate(t *testing.T) {
	tests := []struct {
		name    string
		role    Role
		wantErr bool
		errCode string
	}{
		{
			name: "valid role",
			role: Role{
				ID:          "role-1",
				Name:        "Admin",
				Description: "Administrator role",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			role: Role{
				ID:          "role-1",
				Name:        "",
				Description: "Description",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
		{
			name: "whitespace name",
			role: Role{
				ID:          "role-1",
				Name:        "   ",
				Description: "Description",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
		{
			name: "name too short",
			role: Role{
				ID:          "role-1",
				Name:        "A",
				Description: "Description",
			},
			wantErr: true,
			errCode: errors.CodeNameTooShort,
		},
		{
			name: "name exactly 2 chars",
			role: Role{
				ID:          "role-1",
				Name:        "AB",
				Description: "Description",
			},
			wantErr: false,
		},
		{
			name: "name too long",
			role: Role{
				ID:          "role-1",
				Name:        "A" + string(make([]byte, 50)),
				Description: "Description",
			},
			wantErr: true,
			errCode: errors.CodeNameTooLong,
		},
		{
			name: "name exactly 50 chars",
			role: Role{
				ID:          "role-1",
				Name:        string(make([]byte, 50)),
				Description: "Description",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.role.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("Role.Validate() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("Role.Validate() error type = %T, want *errors.AppError", err)
				}
			}
		})
	}
}

func TestRole_Create(t *testing.T) {
	role := &Role{
		ID:          "role-1",
		Name:        "Admin",
		Description: "Administrator role",
	}

	role.Create()

	events := role.Events()
	if len(events) != 1 {
		t.Fatalf("Role.Create() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(RoleCreatedEvent)
	if !ok {
		t.Fatalf("Role.Create() event type = %T, want RoleCreatedEvent", events[0])
	}

	if event.RoleID != role.ID {
		t.Errorf("RoleCreatedEvent.RoleID = %v, want %v", event.RoleID, role.ID)
	}
	if event.Name != role.Name {
		t.Errorf("RoleCreatedEvent.Name = %v, want %v", event.Name, role.Name)
	}
	if event.CreatedAt == "" {
		t.Error("RoleCreatedEvent.CreatedAt should not be empty")
	}
}

func TestRole_Update(t *testing.T) {
	role := &Role{
		ID:          "role-1",
		Name:        "Old Name",
		Description: "Old Description",
	}

	newName := "New Name"
	newDescription := "New Description"
	role.Update(newName, newDescription)

	if role.Name != newName {
		t.Errorf("Role.Name = %v, want %v", role.Name, newName)
	}
	if role.Description != newDescription {
		t.Errorf("Role.Description = %v, want %v", role.Description, newDescription)
	}

	if role.UpdatedAt.IsZero() {
		t.Error("Role.UpdatedAt should be set")
	}

	events := role.Events()
	if len(events) != 1 {
		t.Fatalf("Role.Update() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(RoleUpdatedEvent)
	if !ok {
		t.Fatalf("Role.Update() event type = %T, want RoleUpdatedEvent", events[0])
	}

	if event.RoleID != role.ID {
		t.Errorf("RoleUpdatedEvent.RoleID = %v, want %v", event.RoleID, role.ID)
	}
	if event.OldName != "Old Name" {
		t.Errorf("RoleUpdatedEvent.OldName = %v, want %v", event.OldName, "Old Name")
	}
	if event.NewName != newName {
		t.Errorf("RoleUpdatedEvent.NewName = %v, want %v", event.NewName, newName)
	}
	if event.UpdatedAt == "" {
		t.Error("RoleUpdatedEvent.UpdatedAt should not be empty")
	}
}

func TestRole_ClearEvents(t *testing.T) {
	role := &Role{
		ID:   "role-1",
		Name: "Admin",
	}

	role.Create()
	if len(role.Events()) != 1 {
		t.Fatal("Expected 1 event after Create()")
	}

	role.ClearEvents()
	if len(role.Events()) != 0 {
		t.Errorf("Role.ClearEvents() should clear all events, got %d", len(role.Events()))
	}
}


