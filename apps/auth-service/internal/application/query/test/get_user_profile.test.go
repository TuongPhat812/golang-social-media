package query

import (
	"context"
	"testing"

	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
)

func TestGetUserProfileQuery_Execute(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()

	// Create a test user
	testUser, err := factory.CreateUser("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	if err := repo.Create(*testUser); err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	query := NewGetUserProfileHandler(repo)

	ctx := context.Background()

	t.Run("Successful Get Profile", func(t *testing.T) {
		resp, err := query.Execute(ctx, testUser.ID)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		if resp.ID != testUser.ID {
			t.Errorf("Response.ID = %v, want %v", resp.ID, testUser.ID)
		}
		if resp.Email != testUser.Email {
			t.Errorf("Response.Email = %v, want %v", resp.Email, testUser.Email)
		}
		if resp.Name != testUser.Name {
			t.Errorf("Response.Name = %v, want %v", resp.Name, testUser.Name)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, err := query.Execute(ctx, "non-existent")
		if err == nil {
			t.Error("Execute() should return error for non-existent user")
		}
	})
}

