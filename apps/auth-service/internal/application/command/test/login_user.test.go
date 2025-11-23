package command

import (
	"context"
	"testing"

	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
)

func TestLoginUserHandler_Handle(t *testing.T) {
	// Setup
	repo := memory.NewUserRepository(nil)
	jwtService := jwt.NewService("test-secret", 1, 168)

	// Create a user
	testUser := user.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewLoginUserHandler(repo, jwtService)

	req := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := handler.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if resp.UserID != testUser.ID {
		t.Errorf("LoginResponse.UserID = %v, want %v", resp.UserID, testUser.ID)
	}
	if resp.AccessToken == "" {
		t.Error("LoginResponse.AccessToken should not be empty")
	}
	if resp.RefreshToken == "" {
		t.Error("LoginResponse.RefreshToken should not be empty")
	}
	if resp.ExpiresIn <= 0 {
		t.Error("LoginResponse.ExpiresIn should be greater than 0")
	}
}

func TestLoginUserHandler_Handle_InvalidEmail(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	jwtService := jwt.NewService("test-secret", 1, 168)
	handler := NewLoginUserHandler(repo, jwtService)

	req := auth.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := handler.Handle(context.Background(), req)
	if err == nil {
		t.Error("Handle() should return error for non-existent email")
	}
}

func TestLoginUserHandler_Handle_InvalidPassword(t *testing.T) {
	repo := memory.NewUserRepository(nil)
	jwtService := jwt.NewService("test-secret", 1, 168)

	// Create a user
	testUser := user.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Password: "correctpassword",
		Name:     "Test User",
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewLoginUserHandler(repo, jwtService)

	req := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := handler.Handle(context.Background(), req)
	if err == nil {
		t.Error("Handle() should return error for wrong password")
	}

	if err != memory.ErrInvalidAuth {
		t.Errorf("Handle() error = %v, want %v", err, memory.ErrInvalidAuth)
	}
}

