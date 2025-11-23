package handlers

import (
	"net/http"

	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// TokenHandler handles token-related endpoints
type TokenHandler struct {
	logout      commandcontracts.LogoutUserCommand
	refresh     commandcontracts.RefreshTokenCommand
	revoke      commandcontracts.RevokeTokenCommand
	validate    querycontracts.ValidateTokenQuery
}

// NewTokenHandler creates a new TokenHandler
func NewTokenHandler(
	logout commandcontracts.LogoutUserCommand,
	refresh commandcontracts.RefreshTokenCommand,
	revoke commandcontracts.RevokeTokenCommand,
	validate querycontracts.ValidateTokenQuery,
) *TokenHandler {
	return &TokenHandler{
		logout:   logout,
		refresh:  refresh,
		revoke:   revoke,
		validate: validate,
	}
}

// Mount mounts public token routes to the router group
func (h *TokenHandler) Mount(group *gin.RouterGroup) {
	// Public route
	group.POST("/validate-token", h.validateToken)
}

// MountProtected mounts protected token routes (require JWT middleware)
func (h *TokenHandler) MountProtected(group *gin.RouterGroup) {
	group.POST("/logout", h.logoutUser)
	group.POST("/refresh", h.refreshToken)
	group.POST("/revoke", h.revokeToken)
}

// logoutUser handles POST /auth/logout
func (h *TokenHandler) logoutUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	token := extractTokenFromHeader(c)
	if token == "" {
		c.Error(errors.NewInvalidRequestError("Token is required"))
		return
	}

	err := h.logout.Execute(c.Request.Context(), commandcontracts.LogoutUserCommandRequest{
		UserID: userID.(string),
		Token:  token,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// refreshToken handles POST /auth/refresh
func (h *TokenHandler) refreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	resp, err := h.refresh.Execute(c.Request.Context(), commandcontracts.RefreshTokenCommandRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, auth.RefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	})
}

// revokeToken handles POST /auth/revoke
func (h *TokenHandler) revokeToken(c *gin.Context) {
	var req auth.RevokeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	err := h.revoke.Execute(c.Request.Context(), commandcontracts.RevokeTokenCommandRequest{
		Token: req.Token,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

// validateToken handles POST /auth/validate-token
func (h *TokenHandler) validateToken(c *gin.Context) {
	var req auth.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	resp, err := h.validate.Execute(c.Request.Context(), req.Token)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, auth.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserID: resp.UserID,
	})
}

// extractTokenFromHeader extracts Bearer token from Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	// Extract Bearer token
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

