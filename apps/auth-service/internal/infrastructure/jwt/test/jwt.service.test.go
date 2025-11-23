package jwt

import (
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	secret := "test-secret-key"
	accessExpHours := 1
	refreshExpHours := 168

	service := NewService(secret, accessExpHours, refreshExpHours)

	if service == nil {
		t.Fatal("NewService() should not return nil")
	}

	if len(service.secret) == 0 {
		t.Error("Service.secret should not be empty")
	}

	expectedAccessExp := time.Duration(accessExpHours) * time.Hour
	if service.accessExpiration != expectedAccessExp {
		t.Errorf("Service.accessExpiration = %v, want %v", service.accessExpiration, expectedAccessExp)
	}

	expectedRefreshExp := time.Duration(refreshExpHours) * time.Hour
	if service.refreshExpiration != expectedRefreshExp {
		t.Errorf("Service.refreshExpiration = %v, want %v", service.refreshExpiration, expectedRefreshExp)
	}
}

func TestNewService_DefaultValues(t *testing.T) {
	secret := "test-secret-key"

	// Test with zero/negative values
	service := NewService(secret, 0, 0)

	expectedAccessExp := 1 * time.Hour
	if service.accessExpiration != expectedAccessExp {
		t.Errorf("Service.accessExpiration = %v, want %v (default)", service.accessExpiration, expectedAccessExp)
	}

	expectedRefreshExp := 168 * time.Hour
	if service.refreshExpiration != expectedRefreshExp {
		t.Errorf("Service.refreshExpiration = %v, want %v (default)", service.refreshExpiration, expectedRefreshExp)
	}
}

func TestService_GenerateTokenPair(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	tokenPair, err := service.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("TokenPair.AccessToken should not be empty")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("TokenPair.RefreshToken should not be empty")
	}

	if tokenPair.ExpiresIn <= 0 {
		t.Error("TokenPair.ExpiresIn should be greater than 0")
	}

	// Verify tokens are different
	if tokenPair.AccessToken == tokenPair.RefreshToken {
		t.Error("AccessToken and RefreshToken should be different")
	}
}

func TestService_GenerateTokenPairWithClaims(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"
	roles := []string{"admin", "user"}
	permissions := []string{"read", "write"}

	tokenPair, err := service.GenerateTokenPairWithClaims(userID, roles, permissions)
	if err != nil {
		t.Fatalf("GenerateTokenPairWithClaims() error = %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("TokenPair.AccessToken should not be empty")
	}

	// Validate the access token contains claims
	claims, err := service.ValidateTokenWithClaims(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateTokenWithClaims() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Claims.UserID = %v, want %v", claims.UserID, userID)
	}

	if len(claims.Roles) != len(roles) {
		t.Errorf("Claims.Roles length = %v, want %v", len(claims.Roles), len(roles))
	}

	if len(claims.Permissions) != len(permissions) {
		t.Errorf("Claims.Permissions length = %v, want %v", len(claims.Permissions), len(permissions))
	}
}

func TestService_ValidateToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	tokenPair, err := service.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	validatedUserID, err := service.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateToken() userID = %v, want %v", validatedUserID, userID)
	}
}

func TestService_ValidateToken_InvalidToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)

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
			_, err := service.ValidateToken(tt.token)
			if err == nil {
				t.Error("ValidateToken() should return error for invalid token")
			}
		})
	}
}

func TestService_ValidateTokenWithClaims(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"
	roles := []string{"admin"}
	permissions := []string{"read"}

	tokenPair, err := service.GenerateTokenPairWithClaims(userID, roles, permissions)
	if err != nil {
		t.Fatalf("GenerateTokenPairWithClaims() error = %v", err)
	}

	claims, err := service.ValidateTokenWithClaims(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateTokenWithClaims() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Claims.UserID = %v, want %v", claims.UserID, userID)
	}

	if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
		t.Errorf("Claims.Roles = %v, want %v", claims.Roles, roles)
	}

	if len(claims.Permissions) != 1 || claims.Permissions[0] != "read" {
		t.Errorf("Claims.Permissions = %v, want %v", claims.Permissions, permissions)
	}
}

func TestService_ValidateRefreshToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	tokenPair, err := service.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	validatedUserID, err := service.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateRefreshToken() userID = %v, want %v", validatedUserID, userID)
	}
}

func TestService_ValidateRefreshToken_WithAccessToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	tokenPair, err := service.GenerateTokenPair(userID)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	// Try to validate access token as refresh token (should fail)
	_, err = service.ValidateRefreshToken(tokenPair.AccessToken)
	if err == nil {
		t.Error("ValidateRefreshToken() should return error when validating access token")
	}
}

func TestService_GenerateToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	token, err := service.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() should not return empty token")
	}

	// Validate the token
	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateToken() userID = %v, want %v", validatedUserID, userID)
	}
}

func TestService_GenerateRefreshToken(t *testing.T) {
	service := NewService("test-secret", 1, 168)
	userID := "user-123"

	refreshToken, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if refreshToken == "" {
		t.Error("GenerateRefreshToken() should not return empty token")
	}

	// Validate the refresh token
	validatedUserID, err := service.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateRefreshToken() userID = %v, want %v", validatedUserID, userID)
	}
}

func TestService_DifferentSecrets(t *testing.T) {
	service1 := NewService("secret-1", 1, 168)
	service2 := NewService("secret-2", 1, 168)
	userID := "user-123"

	token1, err := service1.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Token from service1 should not be valid in service2
	_, err = service2.ValidateToken(token1)
	if err == nil {
		t.Error("Token from different secret should not be valid")
	}
}

