package permission

import (
	"testing"

	"golang-social-media/pkg/errors"
)

func TestPermission_Validate(t *testing.T) {
	tests := []struct {
		name    string
		perm    Permission
		wantErr bool
		errCode string
	}{
		{
			name: "valid permission",
			perm: Permission{
				ID:       "perm-1",
				Name:     "Create Chat",
				Resource: "chat",
				Action:   "create",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			perm: Permission{
				ID:       "perm-1",
				Name:     "",
				Resource: "chat",
				Action:   "create",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
		{
			name: "whitespace name",
			perm: Permission{
				ID:       "perm-1",
				Name:     "   ",
				Resource: "chat",
				Action:   "create",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
		{
			name: "empty resource",
			perm: Permission{
				ID:       "perm-1",
				Name:     "Create Chat",
				Resource: "",
				Action:   "create",
			},
			wantErr: true,
			errCode: "resource_required",
		},
		{
			name: "whitespace resource",
			perm: Permission{
				ID:       "perm-1",
				Name:     "Create Chat",
				Resource: "   ",
				Action:   "create",
			},
			wantErr: true,
			errCode: "resource_required",
		},
		{
			name: "empty action",
			perm: Permission{
				ID:       "perm-1",
				Name:     "Create Chat",
				Resource: "chat",
				Action:   "",
			},
			wantErr: true,
			errCode: "action_required",
		},
		{
			name: "whitespace action",
			perm: Permission{
				ID:       "perm-1",
				Name:     "Create Chat",
				Resource: "chat",
				Action:   "   ",
			},
			wantErr: true,
			errCode: "action_required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.perm.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Permission.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("Permission.Validate() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("Permission.Validate() error type = %T, want *errors.AppError", err)
				}
			}
		})
	}
}

func TestPermission_Create(t *testing.T) {
	perm := &Permission{
		ID:       "perm-1",
		Name:     "Create Chat",
		Resource: "chat",
		Action:   "create",
	}

	perm.Create()

	events := perm.Events()
	if len(events) != 1 {
		t.Fatalf("Permission.Create() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(PermissionCreatedEvent)
	if !ok {
		t.Fatalf("Permission.Create() event type = %T, want PermissionCreatedEvent", events[0])
	}

	if event.PermissionID != perm.ID {
		t.Errorf("PermissionCreatedEvent.PermissionID = %v, want %v", event.PermissionID, perm.ID)
	}
	if event.Name != perm.Name {
		t.Errorf("PermissionCreatedEvent.Name = %v, want %v", event.Name, perm.Name)
	}
	if event.Resource != perm.Resource {
		t.Errorf("PermissionCreatedEvent.Resource = %v, want %v", event.Resource, perm.Resource)
	}
	if event.Action != perm.Action {
		t.Errorf("PermissionCreatedEvent.Action = %v, want %v", event.Action, perm.Action)
	}
	if event.CreatedAt == "" {
		t.Error("PermissionCreatedEvent.CreatedAt should not be empty")
	}
}

func TestPermission_ClearEvents(t *testing.T) {
	perm := &Permission{
		ID:       "perm-1",
		Name:     "Create Chat",
		Resource: "chat",
		Action:   "create",
	}

	perm.Create()
	if len(perm.Events()) != 1 {
		t.Fatal("Expected 1 event after Create()")
	}

	perm.ClearEvents()
	if len(perm.Events()) != 0 {
		t.Errorf("Permission.ClearEvents() should clear all events, got %d", len(perm.Events()))
	}
}

