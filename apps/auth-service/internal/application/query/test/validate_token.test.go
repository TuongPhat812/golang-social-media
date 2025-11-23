package query

import (
	"context"
	"errors"
	"testing"

	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
)

// mockTokenBlacklistRepository is a mock implementation for testing
type mockTokenBlacklistRepository struct {
	blacklisted map[string]bool
	isBlacklisted func(ctx context.Context, tokenID string) (bool, error)
}

func newMockTokenBlacklistRepository() *mockTokenBlacklistRepository {
	return &mockTokenBlacklistRepository{
		blacklisted: make(map[string]bool),
	}
}

func (m *mockTokenBlacklistRepository) Add(ctx context.Context, tokenID string, expirationSeconds int64) error {
	m.blacklisted[tokenID] = true
	return nil
}

func (m *mockTokenBlacklistRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	if m.isBlacklisted != nil {
		return m.isBlacklisted(ctx, tokenID)
	}
	return m.blacklisted[tokenID], nil
}

func (m *mockTokenBlacklistRepository) Remove(ctx context.Context, tokenID string) error {
	delete(m.blacklisted, tokenID)
	return nil
}

func TestValidateTokenQuery_Execute(t *testing.T) {
	jwtService := jwt.NewService("test-secret", 1, 168)
	blacklistRepo := newMockTokenBlacklistRepository()

	query := NewValidateTokenQuery(jwtService, blacklistRepo)

	// Generate a valid token
	userID := "user-123"
	tokenPair, err := jwtService.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	resp, err := query.Execute(context.Background(), tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !resp.Valid {
		t.Error("ValidateTokenQueryResponse.Valid should be true for valid token")
	}
	if resp.UserID != userID {
		t.Errorf("ValidateTokenQueryResponse.UserID = %v, want %v", resp.UserID, userID)
	}
}

func TestValidateTokenQuery_Execute_InvalidToken(t *testing.T) {
	jwtService := jwt.NewService("test-secret", 1, 168)
	blacklistRepo := newMockTokenBlacklistRepository()

	query := NewValidateTokenQuery(jwtService, blacklistRepo)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format",
			token: "not.a.valid.token",
		},
		{
			name:  "wrong secret",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := query.Execute(context.Background(), tt.token)
			if err != nil {
				t.Fatalf("Execute() should not return error, got %v", err)
			}
			if resp.Valid {
				t.Error("ValidateTokenQueryResponse.Valid should be false for invalid token")
			}
			if resp.UserID != "" {
				t.Error("ValidateTokenQueryResponse.UserID should be empty for invalid token")
			}
		})
	}
}

func TestValidateTokenQuery_Execute_BlacklistedToken(t *testing.T) {
	jwtService := jwt.NewService("test-secret", 1, 168)
	blacklistRepo := newMockTokenBlacklistRepository()

	query := NewValidateTokenQuery(jwtService, blacklistRepo)

	// Generate a valid token
	userID := "user-123"
	tokenPair, err := jwtService.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Blacklist the token
	tokenID := user.NewTokenID(tokenPair.AccessToken)
	if err := blacklistRepo.Add(context.Background(), tokenID.String(), 3600); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	resp, err := query.Execute(context.Background(), tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if resp.Valid {
		t.Error("ValidateTokenQueryResponse.Valid should be false for blacklisted token")
	}
	if resp.UserID != "" {
		t.Error("ValidateTokenQueryResponse.UserID should be empty for blacklisted token")
	}
}

func TestValidateTokenQuery_Execute_BlacklistCheckError(t *testing.T) {
	jwtService := jwt.NewService("test-secret", 1, 168)
	blacklistRepo := newMockTokenBlacklistRepository()
	blacklistRepo.isBlacklisted = func(ctx context.Context, tokenID string) (bool, error) {
		return false, errors.New("database error")
	}

	query := NewValidateTokenQuery(jwtService, blacklistRepo)

	// Generate a valid token
	userID := "user-123"
	tokenPair, err := jwtService.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Should return valid=false on blacklist check error, not an error
	resp, err := query.Execute(context.Background(), tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Execute() should not return error on blacklist check failure, got %v", err)
	}
	if resp.Valid {
		t.Error("ValidateTokenQueryResponse.Valid should be false when blacklist check fails")
	}
}

