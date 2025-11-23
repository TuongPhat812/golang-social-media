package handlers

import (
	"net/http"

	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	appcommand "golang-social-media/apps/auth-service/internal/application/command"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related endpoints
type AuthHandler struct {
	registerUser commandcontracts.RegisterUserCommand
	loginUser    *appcommand.LoginUserHandler
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	registerUser commandcontracts.RegisterUserCommand,
	loginUser *appcommand.LoginUserHandler,
) *AuthHandler {
	return &AuthHandler{
		registerUser: registerUser,
		loginUser:    loginUser,
	}
}

// Mount mounts auth routes to the router group
func (h *AuthHandler) Mount(group *gin.RouterGroup) {
	group.POST("/register", h.register)
	group.POST("/login", h.login)
}

// Register handles user registration (exported for direct use)
func (h *AuthHandler) Register(c *gin.Context) {
	h.register(c)
}

// Login handles user login (exported for direct use)
func (h *AuthHandler) Login(c *gin.Context) {
	h.login(c)
}

// register handles user registration
func (h *AuthHandler) register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	resp, err := h.registerUser.Execute(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// login handles user login
func (h *AuthHandler) login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	resp, err := h.loginUser.Handle(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

