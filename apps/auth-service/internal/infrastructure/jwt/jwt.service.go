package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret            []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

func NewService(secret string, accessExpirationHours int, refreshExpirationHours int) *Service {
	if accessExpirationHours <= 0 {
		accessExpirationHours = 1 // Default 1 hour for access token
	}
	if refreshExpirationHours <= 0 {
		refreshExpirationHours = 168 // Default 7 days for refresh token
	}
	return &Service{
		secret:            []byte(secret),
		accessExpiration:  time.Duration(accessExpirationHours) * time.Hour,
		refreshExpiration: time.Duration(refreshExpirationHours) * time.Hour,
	}
}

// TokenClaims represents JWT token claims with roles and permissions
type TokenClaims struct {
	UserID      string
	Roles       []string
	Permissions []string
}

// GenerateTokenPair generates both access and refresh tokens
func (s *Service) GenerateTokenPair(userID string) (*TokenPair, error) {
	return s.GenerateTokenPairWithClaims(userID, []string{}, []string{})
}

// GenerateTokenPairWithClaims generates both access and refresh tokens with roles and permissions
func (s *Service) GenerateTokenPairWithClaims(userID string, roles []string, permissions []string) (*TokenPair, error) {
	accessToken, err := s.generateTokenWithClaims(userID, roles, permissions, s.accessExpiration, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, s.refreshExpiration, "refresh")
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiration.Seconds()),
	}, nil
}

// GenerateToken generates a JWT token for a user (backward compatibility)
func (s *Service) GenerateToken(userID string) (string, error) {
	return s.generateToken(userID, s.accessExpiration, "access")
}

// generateToken generates a JWT token with specified expiration and type
func (s *Service) generateToken(userID string, expiration time.Duration, tokenType string) (string, error) {
	return s.generateTokenWithClaims(userID, []string{}, []string{}, expiration, tokenType)
}

// generateTokenWithClaims generates a JWT token with roles and permissions
func (s *Service) generateTokenWithClaims(userID string, roles []string, permissions []string, expiration time.Duration, tokenType string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":     userID,
		"type":        tokenType,
		"exp":         now.Add(expiration).Unix(),
		"iat":         now.Unix(),
		"roles":       roles,
		"permissions": permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token
func (s *Service) GenerateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, s.refreshExpiration, "refresh")
}

// ValidateToken validates a JWT token and returns the user ID
func (s *Service) ValidateToken(tokenString string) (string, error) {
	claims, err := s.ValidateTokenWithClaims(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ValidateTokenWithClaims validates a JWT token and returns full claims
func (s *Service) ValidateTokenWithClaims(tokenString string) (*TokenClaims, error) {
	return s.validateTokenWithClaims(tokenString, "access")
}

// ValidateRefreshToken validates a refresh token and returns the user ID
func (s *Service) ValidateRefreshToken(tokenString string) (string, error) {
	return s.validateToken(tokenString, "refresh")
}

// validateToken validates a JWT token with specified type (backward compatibility)
func (s *Service) validateToken(tokenString string, expectedType string) (string, error) {
	claims, err := s.validateTokenWithClaims(tokenString, expectedType)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// validateTokenWithClaims validates a JWT token with specified type and returns full claims
func (s *Service) validateTokenWithClaims(tokenString string, expectedType string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Validate token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("missing user_id in token")
	}

	// Extract roles
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		roles = make([]string, 0, len(rolesInterface))
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	// Extract permissions
	var permissions []string
	if permsInterface, ok := claims["permissions"].([]interface{}); ok {
		permissions = make([]string, 0, len(permsInterface))
		for _, p := range permsInterface {
			if permStr, ok := p.(string); ok {
				permissions = append(permissions, permStr)
			}
		}
	}

	return &TokenClaims{
		UserID:      userID,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

