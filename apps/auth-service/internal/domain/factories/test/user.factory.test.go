package factories

import (
	"strings"
	"testing"

	"golang-social-media/apps/auth-service/internal/domain/user"
)

func TestNewUserFactory(t *testing.T) {
	factory := NewUserFactory()

	if factory == nil {
		t.Fatal("NewUserFactory() should not return nil")
	}
}

func TestNewUserFactoryWithIDGenerator(t *testing.T) {
	customID := "custom-user-id"
	idGenerator := func() string {
		return customID
	}

	factory := NewUserFactoryWithIDGenerator(idGenerator)

	if factory == nil {
		t.Fatal("NewUserFactoryWithIDGenerator() should not return nil")
	}

	// Test that custom ID generator is used
	user, err := factory.CreateUser("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.ID != customID {
		t.Errorf("User.ID = %v, want %v", user.ID, customID)
	}
}

func TestUserFactoryImpl_CreateUser(t *testing.T) {
	factory := NewUserFactory()

	tests := []struct {
		name     string
		email    string
		password string
		userName string
		wantErr  bool
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			password: "password123",
			userName: "Test User",
			wantErr:  false,
		},
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			userName: "Test User",
			wantErr:  true,
		},
		{
			name:     "empty password",
			email:    "test@example.com",
			password: "",
			userName: "Test User",
			wantErr:  true,
		},
		{
			name:     "empty name",
			email:    "test@example.com",
			password: "password123",
			userName: "",
			wantErr:  true,
		},
		{
			name:     "invalid email format",
			email:    "notanemail",
			password: "password123",
			userName: "Test User",
			wantErr:  true,
		},
		{
			name:     "password too short",
			email:    "test@example.com",
			password: "12345",
			userName: "Test User",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := factory.CreateUser(tt.email, tt.password, tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Fatal("CreateUser() should not return nil user when no error")
				}
				if user.Email != tt.email {
					t.Errorf("User.Email = %v, want %v", user.Email, tt.email)
				}
				if user.Password != tt.password {
					t.Errorf("User.Password = %v, want %v", user.Password, tt.password)
				}
				if user.Name != tt.userName {
					t.Errorf("User.Name = %v, want %v", user.Name, tt.userName)
				}
				if user.ID == "" {
					t.Error("User.ID should not be empty")
				}
				// Check that domain events are created
				events := user.Events()
				if len(events) == 0 {
					t.Error("User should have domain events after creation")
				}
			}
		})
	}
}

func TestUserFactoryImpl_CreateUser_GeneratesUniqueIDs(t *testing.T) {
	factory := NewUserFactory()

	user1, err1 := factory.CreateUser("user1@example.com", "password123", "User 1")
	if err1 != nil {
		t.Fatalf("CreateUser() error = %v", err1)
	}

	user2, err2 := factory.CreateUser("user2@example.com", "password123", "User 2")
	if err2 != nil {
		t.Fatalf("CreateUser() error = %v", err2)
	}

	if user1.ID == user2.ID {
		t.Error("CreateUser() should generate unique IDs")
	}
}

func TestUserFactoryError_Error(t *testing.T) {
	err := &UserFactoryError{
		Message: "test error",
	}

	if err.Error() != "test error" {
		t.Errorf("UserFactoryError.Error() = %v, want %v", err.Error(), "test error")
	}
}

func TestUserFactoryError_Error_WithCause(t *testing.T) {
	cause := &user.User{}
	err := &UserFactoryError{
		Message: "test error",
		Cause:   cause,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("UserFactoryError.Error() should not be empty")
	}
	if !strings.Contains(errorMsg, "test error") {
		t.Errorf("UserFactoryError.Error() should contain message, got %v", errorMsg)
	}
}

func TestUserFactoryError_Unwrap(t *testing.T) {
	cause := &user.User{}
	err := &UserFactoryError{
		Message: "test error",
		Cause:   cause,
	}

	if err.Unwrap() != cause {
		t.Errorf("UserFactoryError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

