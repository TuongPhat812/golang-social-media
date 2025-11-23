package user

import (
	"testing"
	"time"

	"golang-social-media/pkg/errors"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errCode string
	}{
		{
			name: "valid user",
			user: User{
				ID:       "user-1",
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			user: User{
				ID:       "user-1",
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: true,
			errCode: errors.CodeEmailRequired,
		},
		{
			name: "whitespace email",
			user: User{
				ID:       "user-1",
				Email:    "   ",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: true,
			errCode: errors.CodeEmailRequired,
		},
		{
			name: "invalid email format",
			user: User{
				ID:       "user-1",
				Email:    "notanemail",
				Password: "password123",
				Name:     "Test User",
			},
			wantErr: true,
			errCode: errors.CodeEmailInvalid,
		},
		{
			name: "empty password",
			user: User{
				ID:       "user-1",
				Email:    "test@example.com",
				Password: "",
				Name:     "Test User",
			},
			wantErr: true,
			errCode: errors.CodePasswordRequired,
		},
		{
			name: "password too short",
			user: User{
				ID:       "user-1",
				Email:    "test@example.com",
				Password: "12345",
				Name:     "Test User",
			},
			wantErr: true,
			errCode: errors.CodePasswordTooShort,
		},
		{
			name: "empty name",
			user: User{
				ID:       "user-1",
				Email:    "test@example.com",
				Password: "password123",
				Name:     "",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
		{
			name: "whitespace name",
			user: User{
				ID:       "user-1",
				Email:    "test@example.com",
				Password: "password123",
				Name:     "   ",
			},
			wantErr: true,
			errCode: errors.CodeNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("User.Validate() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("User.Validate() error type = %T, want *errors.AppError", err)
				}
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	user := User{}
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errCode  string
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
			errCode:  errors.CodePasswordRequired,
		},
		{
			name:     "whitespace password",
			password: "   ",
			wantErr:  true,
			errCode:  errors.CodePasswordRequired,
		},
		{
			name:     "password too short",
			password: "12345",
			wantErr:  true,
			errCode:  errors.CodePasswordTooShort,
		},
		{
			name:     "password exactly 6 chars",
			password: "123456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("User.ValidatePassword() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("User.ValidatePassword() error type = %T, want *errors.AppError", err)
				}
			}
		})
	}
}

func TestUser_Create(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	user.Create()

	events := user.Events()
	if len(events) != 1 {
		t.Fatalf("User.Create() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(UserCreatedEvent)
	if !ok {
		t.Fatalf("User.Create() event type = %T, want UserCreatedEvent", events[0])
	}

	if event.UserID != user.ID {
		t.Errorf("UserCreatedEvent.UserID = %v, want %v", event.UserID, user.ID)
	}
	if event.Email != user.Email {
		t.Errorf("UserCreatedEvent.Email = %v, want %v", event.Email, user.Email)
	}
	if event.Name != user.Name {
		t.Errorf("UserCreatedEvent.Name = %v, want %v", event.Name, user.Name)
	}
	if event.CreatedAt == "" {
		t.Error("UserCreatedEvent.CreatedAt should not be empty")
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Old Name",
	}

	newName := "New Name"
	user.UpdateProfile(newName)

	if user.Name != newName {
		t.Errorf("User.Name = %v, want %v", user.Name, newName)
	}

	if user.UpdatedAt.IsZero() {
		t.Error("User.UpdatedAt should be set")
	}

	events := user.Events()
	if len(events) != 1 {
		t.Fatalf("User.UpdateProfile() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(UserProfileUpdatedEvent)
	if !ok {
		t.Fatalf("User.UpdateProfile() event type = %T, want UserProfileUpdatedEvent", events[0])
	}

	if event.UserID != user.ID {
		t.Errorf("UserProfileUpdatedEvent.UserID = %v, want %v", event.UserID, user.ID)
	}
	if event.OldName != "Old Name" {
		t.Errorf("UserProfileUpdatedEvent.OldName = %v, want %v", event.OldName, "Old Name")
	}
	if event.NewName != newName {
		t.Errorf("UserProfileUpdatedEvent.NewName = %v, want %v", event.NewName, newName)
	}
}

func TestUser_UpdateProfile_EmptyName(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Old Name",
	}

	user.UpdateProfile("")

	// Should not update if name is empty
	if user.Name != "Old Name" {
		t.Errorf("User.Name = %v, want %v", user.Name, "Old Name")
	}

	events := user.Events()
	if len(events) != 0 {
		t.Errorf("User.UpdateProfile() with empty name should not add event, got %d", len(events))
	}
}

func TestUser_ChangePassword(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "oldpassword",
		Name:     "Test User",
	}

	newPassword := "newpassword123"
	user.ChangePassword(newPassword)

	if user.Password != newPassword {
		t.Errorf("User.Password = %v, want %v", user.Password, newPassword)
	}

	if user.UpdatedAt.IsZero() {
		t.Error("User.UpdatedAt should be set")
	}

	events := user.Events()
	if len(events) != 1 {
		t.Fatalf("User.ChangePassword() should add 1 event, got %d", len(events))
	}

	event, ok := events[0].(UserPasswordChangedEvent)
	if !ok {
		t.Fatalf("User.ChangePassword() event type = %T, want UserPasswordChangedEvent", events[0])
	}

	if event.UserID != user.ID {
		t.Errorf("UserPasswordChangedEvent.UserID = %v, want %v", event.UserID, user.ID)
	}
	if event.UpdatedAt == "" {
		t.Error("UserPasswordChangedEvent.UpdatedAt should not be empty")
	}
}

func TestUser_ClearEvents(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	user.Create()
	if len(user.Events()) != 1 {
		t.Fatal("Expected 1 event after Create()")
	}

	user.ClearEvents()
	if len(user.Events()) != 0 {
		t.Errorf("User.ClearEvents() should clear all events, got %d", len(user.Events()))
	}
}

func TestUser_Events(t *testing.T) {
	user := &User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Initially no events
	if len(user.Events()) != 0 {
		t.Errorf("New User should have 0 events, got %d", len(user.Events()))
	}

	// Add multiple events
	user.Create()
	user.UpdateProfile("New Name")
	user.ChangePassword("newpass")

	events := user.Events()
	if len(events) != 3 {
		t.Errorf("User should have 3 events, got %d", len(events))
	}
}

