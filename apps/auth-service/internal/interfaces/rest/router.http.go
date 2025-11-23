package rest

import (
	"net/http"
	"os"
	"time"

	appcommand "golang-social-media/apps/auth-service/internal/application/command"
	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/interfaces/rest/handlers"
	"golang-social-media/apps/auth-service/internal/interfaces/rest/middleware"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/pkg/cache"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Auth     *handlers.AuthHandler
	Profile  *handlers.ProfileHandler
	Password *handlers.PasswordHandler
	Token    *handlers.TokenHandler
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	registerUser commandcontracts.RegisterUserCommand,
	loginUser *appcommand.LoginUserHandler,
) *handlers.AuthHandler {
	return handlers.NewAuthHandler(registerUser, loginUser)
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(
	updateProfile commandcontracts.UpdateProfileCommand,
	getProfile querycontracts.GetUserProfileQuery,
	getCurrentUser querycontracts.GetCurrentUserQuery,
) *handlers.ProfileHandler {
	return handlers.NewProfileHandler(updateProfile, getProfile, getCurrentUser)
}

// NewPasswordHandler creates a new PasswordHandler
func NewPasswordHandler(changePassword commandcontracts.ChangePasswordCommand) *handlers.PasswordHandler {
	return handlers.NewPasswordHandler(changePassword)
}

// NewTokenHandler creates a new TokenHandler
func NewTokenHandler(
	logout commandcontracts.LogoutUserCommand,
	refresh commandcontracts.RefreshTokenCommand,
	revoke commandcontracts.RevokeTokenCommand,
	validate querycontracts.ValidateTokenQuery,
) *handlers.TokenHandler {
	return handlers.NewTokenHandler(logout, refresh, revoke, validate)
}

// NewHandlers creates all HTTP handlers
func NewHandlers(
	authHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	passwordHandler *handlers.PasswordHandler,
	tokenHandler *handlers.TokenHandler,
) *Handlers {
	return &Handlers{
		Auth:     authHandler,
		Profile:  profileHandler,
		Password: passwordHandler,
		Token:    tokenHandler,
	}
}

// NewRouter creates and configures the HTTP router
func NewRouter(h *Handlers, jwtService *jwt.Service, cache cache.Cache) *gin.Engine {
	router := gin.New()

	// Initialize error transformer
	devMode := os.Getenv("ENV") == "development"
	transformer := errors.NewTransformer(devMode)

	// Global middlewares (applied to all routes)
	// Order matters: RequestID -> IP Filter -> Size Limiter -> Compression -> Cache Control -> Metrics -> Request Log -> Error -> Recovery
	router.Use(
		middleware.RequestIDMiddleware(),                    // Request ID tracking (first)
		middleware.IPFilterMiddleware(middleware.DefaultIPFilterConfig()), // IP whitelist/blacklist
		middleware.SizeLimiterMiddleware(middleware.DefaultSizeLimiterConfig()), // Request size limit
		middleware.CompressionMiddleware(),                  // Gzip compression (based on Accept-Encoding)
		middleware.CacheControlMiddleware(middleware.DefaultCacheControlConfig()), // Cache control headers
		middleware.MetricsMiddleware(nil),                  // Metrics collection (simple in-memory)
		middleware.RequestLogMiddleware(),                  // Structured request logging
		gin.Recovery(),                                      // Panic recovery
		errors.ErrorMiddleware(transformer),                 // Error handling
		middleware.CORSMiddleware(middleware.DefaultCORSConfig()), // CORS
		middleware.TimeoutMiddleware(30*time.Second),        // Request timeout
	)

	// Health check (no rate limiting)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "auth-service OK"})
	})

	// Auth routes group
	authGroup := router.Group("/auth")

	// General rate limiter: 100 requests per minute per IP
	authGroup.Use(middleware.RateLimiterMiddleware(cache, middleware.RateLimiterConfig{
		Requests:     100,
		Window:       1 * time.Minute,
		KeyFunc:      middleware.GetClientIP,
		ErrorMessage: "Too many requests. Please try again later.",
	}))

	// Stricter rate limiter for login/register endpoints
	loginRegisterGroup := authGroup.Group("")
	loginRegisterGroup.Use(middleware.RateLimiterMiddleware(cache, middleware.RateLimiterConfig{
		Requests:     5,  // 5 requests
		Window:       1 * time.Minute,  // per minute
		KeyFunc:      middleware.GetClientIP,
		SkipFunc: func(c *gin.Context) bool {
			// Only apply to login/register endpoints
			return c.Request.URL.Path != "/auth/login" && c.Request.URL.Path != "/auth/register"
		},
		ErrorMessage: "Too many login/register attempts. Please try again in a minute.",
	}))
	// Mount login/register with stricter rate limiting
	loginRegisterGroup.POST("/login", h.Auth.Login)
	loginRegisterGroup.POST("/register", h.Auth.Register)

	// Other public routes (with general rate limiting only)
	h.Profile.Mount(authGroup)
	h.Token.Mount(authGroup)

	// Protected routes (require JWT)
	protected := authGroup.Group("")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		// Profile protected routes
		h.Profile.MountProtected(protected)

		// Password protected routes
		h.Password.Mount(protected)

		// Token protected routes
		h.Token.MountProtected(protected)
	}

	return router
}
