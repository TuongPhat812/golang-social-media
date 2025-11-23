package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	pkgerrors "golang-social-media/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTokenBlacklistRepository is a mock implementation for testing
type MockTokenBlacklistRepository struct {
	mock.Mock
}

func (m *MockTokenBlacklistRepository) AddToken(ctx context.Context, tokenID string, expiration time.Duration) error {
	args := m.Called(ctx, tokenID, expiration)
	return args.Error(0)
}

func (m *MockTokenBlacklistRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	args := m.Called(ctx, tokenID)
	return args.Bool(0), args.Error(1)
}

func TestRefreshTokenCommand_Execute(t *testing.T) {
	ctx := context.Background()
	userRepo := memory.NewUserRepository(nil)
	jwtService := jwt.NewService("test-secret", 1, 168)
	mockBlacklistRepo := new(MockTokenBlacklistRepository)

	// Create a test user
	factory := factories.NewUserFactory()
	testUser, err := factory.CreateUser("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	if err := userRepo.Create(*testUser); err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	// Generate a valid refresh token
	tokenPair, err := jwtService.GenerateTokenPair(testUser.ID)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	t.Run("Successful Token Refresh", func(t *testing.T) {
		mockBlacklistRepo.On("IsBlacklisted", ctx, mock.AnythingOfType("string")).Return(false, nil)
		mockBlacklistRepo.On("AddToken", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

		cmd := NewRefreshTokenCommand(userRepo, jwtService, mockBlacklistRepo)
		req := contracts.RefreshTokenCommandRequest{
			RefreshToken: tokenPair.RefreshToken,
		}

		resp, err := cmd.Execute(ctx, req)

		assert.Nil(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Greater(t, resp.ExpiresIn, int64(0))

		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		cmd := NewRefreshTokenCommand(userRepo, jwtService, mockBlacklistRepo)
		req := contracts.RefreshTokenCommandRequest{
			RefreshToken: "invalid-token",
		}

		_, err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.IsType(t, &pkgerrors.AppError{}, err)
		mockBlacklistRepo.AssertNotCalled(t, "IsBlacklisted", mock.Anything, mock.Anything)
	})

	t.Run("Blacklisted Refresh Token", func(t *testing.T) {
		mockBlacklistRepo.On("IsBlacklisted", ctx, mock.AnythingOfType("string")).Return(true, nil)

		cmd := NewRefreshTokenCommand(userRepo, jwtService, mockBlacklistRepo)
		req := contracts.RefreshTokenCommandRequest{
			RefreshToken: tokenPair.RefreshToken,
		}

		_, err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.IsType(t, &pkgerrors.AppError{}, err)
		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Create a token for a non-existent user
		nonExistentUserID := "non-existent-user"
		tokenPair2, err := jwtService.GenerateTokenPair(nonExistentUserID)
		if err != nil {
			t.Fatalf("Failed to generate token pair: %v", err)
		}

		mockBlacklistRepo.On("IsBlacklisted", ctx, mock.AnythingOfType("string")).Return(false, nil)

		cmd := NewRefreshTokenCommand(userRepo, jwtService, mockBlacklistRepo)
		req := contracts.RefreshTokenCommandRequest{
			RefreshToken: tokenPair2.RefreshToken,
		}

		_, err = cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.IsType(t, &pkgerrors.AppError{}, err)
		mockBlacklistRepo.AssertExpectations(t)
	})

	t.Run("Blacklist Check Error", func(t *testing.T) {
		blacklistErr := errors.New("redis connection error")
		mockBlacklistRepo.On("IsBlacklisted", ctx, mock.AnythingOfType("string")).Return(false, blacklistErr)

		cmd := NewRefreshTokenCommand(userRepo, jwtService, mockBlacklistRepo)
		req := contracts.RefreshTokenCommandRequest{
			RefreshToken: tokenPair.RefreshToken,
		}

		_, err := cmd.Execute(ctx, req)

		assert.NotNil(t, err)
		assert.Equal(t, blacklistErr, err)
		mockBlacklistRepo.AssertExpectations(t)
	})
}

