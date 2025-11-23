package command

import (
	"context"
	"errors"
	"testing"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
)

// mockEventDispatcher is a mock implementation of EventDispatcher for testing
type mockUpdateProfileEventDispatcher struct {
	dispatch func(ctx context.Context, event interface{}) error
}

func newMockUpdateProfileEventDispatcher() *mockUpdateProfileEventDispatcher {
	return &mockUpdateProfileEventDispatcher{}
}

func (m *mockUpdateProfileEventDispatcher) Dispatch(ctx context.Context, event interface{}) error {
	if m.dispatch != nil {
		return m.dispatch(ctx, event)
	}
	return nil
}

func TestUpdateProfileCommand_Execute(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockUpdateProfileEventDispatcher()

	// Create a test user
	testUser, err := factory.CreateUser("test@example.com", "password123", "Old Name")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	if err := repo.Create(*testUser); err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	cmd := NewUpdateProfileCommand(repo, factory, dispatcher)

	ctx := context.Background()
	req := contracts.UpdateProfileCommandRequest{
		UserID: testUser.ID,
		Name:   "New Name",
	}

	t.Run("Successful Update", func(t *testing.T) {
		resp, err := cmd.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		if resp.ID != testUser.ID {
			t.Errorf("Response.ID = %v, want %v", resp.ID, testUser.ID)
		}
		if resp.Name != "New Name" {
			t.Errorf("Response.Name = %v, want %v", resp.Name, "New Name")
		}

		// Verify user was updated in repository
		updatedUser, err := repo.GetByID(testUser.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if updatedUser.Name != "New Name" {
			t.Errorf("Updated user name = %v, want %v", updatedUser.Name, "New Name")
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		req := contracts.UpdateProfileCommandRequest{
			UserID: "non-existent",
			Name:   "New Name",
		}

		_, err := cmd.Execute(ctx, req)
		if err == nil {
			t.Error("Execute() should return error for non-existent user")
		}
	})

	t.Run("Empty Name (No Change)", func(t *testing.T) {
		// Reset user name
		testUser.Name = "Old Name"
		repo.Update(*testUser)

		req := contracts.UpdateProfileCommandRequest{
			UserID: testUser.ID,
			Name:   "",
		}

		resp, err := cmd.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		if resp.Name != "Old Name" {
			t.Errorf("Response.Name = %v, want %v (should not change)", resp.Name, "Old Name")
		}
	})

	t.Run("Event Dispatch Error (Should Not Fail)", func(t *testing.T) {
		dispatcher.dispatch = func(ctx context.Context, event interface{}) error {
			return errors.New("dispatch error")
		}

		req := contracts.UpdateProfileCommandRequest{
			UserID: testUser.ID,
			Name:   "Another Name",
		}

		// Command should still succeed even if event dispatch fails
		resp, err := cmd.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Execute() should not fail on event dispatch error, got %v", err)
		}

		if resp.Name != "Another Name" {
			t.Errorf("Response.Name = %v, want %v", resp.Name, "Another Name")
		}
	})
}

