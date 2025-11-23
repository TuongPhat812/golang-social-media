package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheControlConfig configures cache control behavior
type CacheControlConfig struct {
	// MaxAge in seconds (default: 0 = no cache)
	MaxAge int
	// SharedMaxAge in seconds (for CDN/proxy)
	SharedMaxAge int
	// MustRevalidate forces revalidation
	MustRevalidate bool
	// NoCache prevents caching
	NoCache bool
	// NoStore prevents storage
	NoStore bool
	// Private prevents shared caches
	Private bool
	// Public allows caching
	Public bool
}

// DefaultCacheControlConfig returns default cache control config
func DefaultCacheControlConfig() CacheControlConfig {
	return CacheControlConfig{
		MaxAge:         0,
		MustRevalidate: true,
		NoCache:        true,
		NoStore:        false,
		Private:        true,
	}
}

// CacheControlMiddleware creates a cache control middleware
func CacheControlMiddleware(config CacheControlConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build Cache-Control header
		var directives []string

		if config.NoStore {
			directives = append(directives, "no-store")
		}

		if config.NoCache {
			directives = append(directives, "no-cache")
		}

		if config.Private {
			directives = append(directives, "private")
		}

		if config.Public {
			directives = append(directives, "public")
		}

		if config.MaxAge > 0 {
			directives = append(directives, fmt.Sprintf("max-age=%d", config.MaxAge))
		}

		if config.SharedMaxAge > 0 {
			directives = append(directives, fmt.Sprintf("s-maxage=%d", config.SharedMaxAge))
		}

		if config.MustRevalidate {
			directives = append(directives, "must-revalidate")
		}

		if len(directives) > 0 {
			cacheControl := directives[0]
			for i := 1; i < len(directives); i++ {
				cacheControl += ", " + directives[i]
			}
			c.Header("Cache-Control", cacheControl)
		}

		// Set ETag (simple implementation)
		etag := fmt.Sprintf(`"%x"`, c.Request.URL.Path)
		c.Header("ETag", etag)

		// Check If-None-Match header for conditional requests
		if match := c.GetHeader("If-None-Match"); match == etag {
			c.Status(304) // Not Modified
			c.Abort()
			return
		}

		// Set Last-Modified
		c.Header("Last-Modified", time.Now().UTC().Format(time.RFC1123))

		c.Next()
	}
}

