package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang-social-media/pkg/cache"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig configures rate limiter behavior
type RateLimiterConfig struct {
	// Requests per window
	Requests int
	// Time window duration
	Window time.Duration
	// Key function to identify the client (default: by IP)
	KeyFunc func(*gin.Context) string
	// Skip function to bypass rate limiting
	SkipFunc func(*gin.Context) bool
	// Error message when rate limit exceeded
	ErrorMessage string
}

// DefaultRateLimiterConfig returns default rate limiter config
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Requests:     100,              // 100 requests
		Window:       1 * time.Minute,  // per minute
		KeyFunc:      GetClientIP,      // by IP address
		SkipFunc:     nil,              // no skip
		ErrorMessage: "Rate limit exceeded. Please try again later.",
	}
}

// RateLimiterMiddleware creates a rate limiter middleware
// Uses Redis if available, falls back to in-memory limiter
func RateLimiterMiddleware(cache cache.Cache, config RateLimiterConfig) gin.HandlerFunc {
	// If no cache, use in-memory limiter
	if cache == nil {
		logger.Component("auth.middleware.rate_limiter").
			Warn().
			Msg("cache not available, using in-memory rate limiter")
		return inMemoryRateLimiter(config)
	}

	return func(c *gin.Context) {
		// Skip if SkipFunc returns true
		if config.SkipFunc != nil && config.SkipFunc(c) {
			c.Next()
			return
		}

		// Get client identifier
		key := config.KeyFunc(c)
		rateLimitKey := fmt.Sprintf("auth:rate_limit:%s", key)

		ctx := c.Request.Context()

		// Get current count
		countData, err := cache.Get(ctx, rateLimitKey)
		var count int
		if err != nil {
			// Key doesn't exist or error - start fresh
			count = 0
		} else {
			count, _ = strconv.Atoi(string(countData))
		}

		// Check if limit exceeded
		if count >= config.Requests {
			logger.Component("auth.middleware.rate_limiter").
				Warn().
				Str("key", key).
				Int("count", count).
				Int("limit", config.Requests).
				Msg("rate limit exceeded")

			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))
			c.Error(errors.NewTooManyRequestsError(config.ErrorMessage))
			c.Abort()
			return
		}

		// Increment count
		count++
		countStr := strconv.Itoa(count)

		// Set with expiration (window duration)
		if err := cache.Set(ctx, rateLimitKey, []byte(countStr), config.Window); err != nil {
			logger.Component("auth.middleware.rate_limiter").
				Error().
				Err(err).
				Str("key", key).
				Msg("failed to update rate limit counter")
			// Continue on error (fail open)
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(config.Requests-count))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))

		c.Next()
	}
}

// inMemoryRateLimiter creates an in-memory rate limiter (fallback)
func inMemoryRateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	type limitInfo struct {
		count     int
		resetTime time.Time
	}

	limits := make(map[string]*limitInfo)
	var mu sync.RWMutex

	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for key, info := range limits {
				if now.After(info.resetTime) {
					delete(limits, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		if config.SkipFunc != nil && config.SkipFunc(c) {
			c.Next()
			return
		}

		key := config.KeyFunc(c)
		now := time.Now()

		mu.Lock()
		info, exists := limits[key]
		if !exists || now.After(info.resetTime) {
			info = &limitInfo{
				count:     0,
				resetTime: now.Add(config.Window),
			}
			limits[key] = info
		}

		if info.count >= config.Requests {
			mu.Unlock()
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(info.resetTime.Unix(), 10))
			c.Error(errors.NewTooManyRequestsError(config.ErrorMessage))
			c.Abort()
			return
		}

		info.count++
		remaining := config.Requests - info.count
		resetTime := info.resetTime
		mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}

// GetClientIP extracts client IP from request (exported for use in configs)
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := 0; idx < len(ip); idx++ {
			if ip[idx] == ',' {
				ip = ip[:idx]
				break
			}
		}
		return ip
	}

	// Check X-Real-IP header
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}

