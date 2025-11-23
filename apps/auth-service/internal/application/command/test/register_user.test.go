package command

import (
	"context"
	"errors"
	"testing"

	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
	pkgerrors "golang-social-media/pkg/errors"
)

// mockEventDispatcher is a mock implementation of EventDispatcher for testing
type mockEventDispatcher struct {
	dispatch func(ctx context.Context, event user.DomainEvent) error
}

func newMockEventDispatcher() *mockEventDispatcher {
	return &mockEventDispatcher{}
}

func (m *mockEventDispatcher) RegisterHandler(eventType string, handler event_dispatcher.EventHandler) {
	// No-op for testing
}

func (m *mockEventDispatcher) Dispatch(ctx context.Context, event user.DomainEvent) error {
	if m.dispatch != nil {
		return m.dispatch(ctx, event)
	}
	return nil
}

func TestRegisterUserCommand_Execute(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockEventDispatcher()

	cmd := NewRegisterUserCommand(repo, factory, dispatcher)

	req := auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	resp, err := cmd.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if resp.ID == "" {
		t.Error("RegisterResponse.ID should not be empty")
	}
	if resp.Email != req.Email {
		t.Errorf("RegisterResponse.Email = %v, want %v", resp.Email, req.Email)
	}
	if resp.Name != req.Name {
		t.Errorf("RegisterResponse.Name = %v, want %v", resp.Name, req.Name)
	}

	// Verify user was created in repository
	createdUser, err := repo.GetByID(resp.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if createdUser.Email != req.Email {
		t.Errorf("Created user email = %v, want %v", createdUser.Email, req.Email)
	}
}

func TestRegisterUserCommand_Execute_DuplicateEmail(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockEventDispatcher()

	cmd := NewRegisterUserCommand(repo, factory, dispatcher)

	req := auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Create first user
	_, err := cmd.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("First Execute() error = %v", err)
	}

	// Try to create duplicate
	_, err = cmd.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("Execute() should return error for duplicate email")
	}

	if appErr, ok := err.(*pkgerrors.AppError); ok {
		if appErr.Code != pkgerrors.CodeEmailAlreadyExists {
			t.Errorf("Error code = %v, want %v", appErr.Code, pkgerrors.CodeEmailAlreadyExists)
		}
	} else {
		t.Errorf("Error type = %T, want *pkgerrors.AppError", err)
	}
}

func TestRegisterUserCommand_Execute_InvalidUser(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockEventDispatcher()

	cmd := NewRegisterUserCommand(repo, factory, dispatcher)

	tests := []struct {
		name string
		req  auth.RegisterRequest
	}{
		{
			name: "invalid email",
			req: auth.RegisterRequest{
				Email:    "notanemail",
				Password: "password123",
				Name:     "Test User",
			},
		},
		{
			name: "password too short",
			req: auth.RegisterRequest{
				Email:    "test@example.com",
				Password: "12345",
				Name:     "Test User",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.Execute(context.Background(), tt.req)
			if err == nil {
				t.Error("Execute() should return error for invalid user")
			}
		})
	}
}

func TestRegisterUserCommand_Execute_EventDispatchError(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	factory := factories.NewUserFactory()
	dispatcher := newMockEventDispatcher()
	dispatcher.dispatch = func(ctx context.Context, event user.DomainEvent) error {
		return errors.New("event dispatch error")
	}

	cmd := NewRegisterUserCommand(repo, factory, dispatcher)

	req := auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Event dispatch error should not fail the command (it's logged but doesn't fail)
	resp, err := cmd.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute() should not fail on event dispatch error, got %v", err)
	}

	if resp.ID == "" {
		t.Error("RegisterResponse.ID should not be empty even if event dispatch fails")
	}
}

