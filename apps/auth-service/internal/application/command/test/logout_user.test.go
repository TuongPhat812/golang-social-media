package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogoutTokenBlacklistRepository is a mock implementation for testing
type MockLogoutTokenBlacklistRepository struct {
	mock.Mock
}

func (m *MockLogoutTokenBlacklistRepository) AddToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	args := m.Called(ctx, tokenID, expiration)
	return args.Error(0)
}

func (m *MockLogoutTokenBlacklistRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	args := m.Called(ctx, tokenID)
	return args.Bool(0), args.Error(1)
}

func TestLogoutUserCommand_Execute(t *testing.T) {
	ctx := context.Background()
	mockBlacklistRepo := new(MockLogoutTokenBlacklistRepository)

	testToken := "test-token-string"
	testUserID := "user-123"

	t.Run("Successful Logout", func(t *testing.T) {
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

		cmd := NewLogoutUserCommand(mockBlacklistRepo)
		req := contracts.LogoutUserCommandRequest{
			UserID: testUserID,
			Token:   testToken,
		}

		err := cmd.Execute(ctx, req)

		assert.Nil(t, err)
		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("Blacklist Add Fails", func(t *testing.T) {
		blacklistErr := errors.New("redis connection error")
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(blacklistErr)

		cmd := NewLogoutUserCommand(mockBlacklistRepo)
		req := contracts.LogoutUserCommandRequest{
			UserID: testUserID,
			Token:   testToken,
		}

		err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.Equal(t, blacklistErr, err)
		mockBlacklistRepo.AssertExpectations(t)
	})
}

