package command

import (
	"context"
	"errors"
	"testing"

	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	pkgerrors "golang-social-media/pkg/errors"
)

// mockEventDispatcher is a mock implementation of EventDispatcher for testing
type mockChangePasswordEventDispatcher struct {
	dispatch func(ctx context.Context, event interface{}) error
}

func newMockChangePasswordEventDispatcher() *mockChangePasswordEventDispatcher {
	return &mockChangePasswordEventDispatcher{}
}

func (m *mockChangePasswordEventDispatcher) Dispatch(ctx context.Context, event interface{}) error {
	if m.dispatch != nil {
		return m.dispatch(ctx, event)
	}
	return nil
}

func TestChangePasswordCommand_Execute(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockChangePasswordEventDispatcher()

	// Create a test user
	testUser, err := factory.CreateUser("test@example.com", "oldpassword", "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	if err := repo.Create(*testUser); err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	cmd := NewChangePasswordCommand(repo, dispatcher)

	ctx := context.Background()

	t.Run("Successful Password Change", func(t *testing.T) {
		req := contracts.ChangePasswordCommandRequest{
			UserID:          testUser.ID,
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword123",
		}

		err := cmd.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		// Verify password was changed in repository
		updatedUser, err := repo.GetByID(testUser.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if updatedUser.Password != "newpassword123" {
			t.Errorf("Updated user password = %v, want %v", updatedUser.Password, "newpassword123")
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		req := contracts.ChangePasswordCommandRequest{
			UserID:          "non-existent",
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword123",
		}

		err := cmd.Execute(ctx, req)
		if err == nil {
			t.Error("Execute() should return error for non-existent user")
		}
	})

	t.Run("Invalid Current Password", func(t *testing.T) {
		req := contracts.ChangePasswordCommandRequest{
			UserID:          testUser.ID,
			CurrentPassword: "wrongpassword",
			NewPassword:     "newpassword123",
		}

		err := cmd.Execute(ctx, req)
		if err == nil {
			t.Error("Execute() should return error for invalid current password")
		}

		if appErr, ok := err.(*pkgerrors.AppError); ok {
			if appErr.Code != pkgerrors.CodeInvalidCredentials {
				t.Errorf("Error code = %v, want %v", appErr.Code, pkgerrors.CodeInvalidCredentials)
			}
		} else {
			t.Errorf("Error type = %T, want *pkgerrors.AppError", err)
		}
	})

	t.Run("New Password Too Short", func(t *testing.T) {
		req := contracts.ChangePasswordCommandRequest{
			UserID:          testUser.ID,
			CurrentPassword: "oldpassword",
			NewPassword:     "12345", // Too short
		}

		err := cmd.Execute(ctx, req)
		if err == nil {
			t.Error("Execute() should return error for password too short")
		}

		if appErr, ok := err.(*pkgerrors.AppError); ok {
			if appErr.Code != pkgerrors.CodePasswordTooShort {
				t.Errorf("Error code = %v, want %v", appErr.Code, pkgerrors.CodePasswordTooShort)
			}
		}
	})

	t.Run("Event Dispatch Error (Should Not Fail)", func(t *testing.T) {
		dispatcher.dispatch = func(ctx context.Context, event interface{}) error {
			return errors.New("dispatch error")
		}

		req := contracts.ChangePasswordCommandRequest{
			UserID:          testUser.ID,
			CurrentPassword: "oldpassword",
			NewPassword:     "anotherpassword123",
		}

		// Command should still succeed even if event dispatch fails
		err := cmd.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Execute() should not fail on event dispatch error, got %v", err)
		}
	})
}

