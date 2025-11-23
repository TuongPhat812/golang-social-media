package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	pkgerrors "golang-social-media/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTokenBlacklistRepository is a mock implementation for testing
type MockRevokeTokenBlacklistRepository struct {
	mock.Mock
}

func (m *MockRevokeTokenBlacklistRepository) AddToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	args := m.Called(ctx, tokenID, expiration)
	return args.Error(0)
}

func (m *MockRevokeTokenBlacklistRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	args := m.Called(ctx, tokenID)
	return args.Bool(0), args.Error(1)
}

func TestRevokeTokenCommand_Execute(t *testing.T) {
	ctx := context.Background()
	jwtService := jwt.NewService("test-secret", 1, 168)
	mockBlacklistRepo := new(MockRevokeTokenBlacklistRepository)

	testUserID := "user-123"
	tokenPair, err := jwtService.GenerateTokenPair(testUserID)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	t.Run("Successful Access Token Revocation", func(t *testing.T) {
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

		cmd := NewRevokeTokenCommand(jwtService, mockBlacklistRepo)
		req := contracts.RevokeTokenCommandRequest{
			Token: tokenPair.AccessToken,
		}

		err := cmd.Execute(ctx, req)

		assert.Nil(t, err)
		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("Successful Refresh Token Revocation", func(t *testing.T) {
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

		cmd := NewRevokeTokenCommand(jwtService, mockBlacklistRepo)
		req := contracts.RevokeTokenCommandRequest{
			Token: tokenPair.RefreshToken,
		}

		err := cmd.Execute(ctx, req)

		assert.Nil(t, err)
		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		cmd := NewRevokeTokenCommand(jwtService, mockBlacklistRepo)
		req := contracts.RevokeTokenCommandRequest{
			Token: "invalid-token",
		}

		err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.IsType(t, &pkgerrors.AppError{}, err)
		mockBlacklistRepo.AssertNotCalled(t, "AddToken", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Blacklist Add Fails", func(t *testing.T) {
		blacklistErr := errors.New("redis connection error")
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(blacklistErr)

		cmd := NewRevokeTokenCommand(jwtService, mockBlacklistRepo)
		req := contracts.RevokeTokenCommandRequest{
			Token: tokenPair.AccessToken,
		}

		err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.Equal(t, blacklistErr, err)
		mockBlacklistRepo.AssertExpectations(t)
	})
}

