package middleware

import (
	"net"
	"net/http"
	"strings"

	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// IPFilterConfig configures IP filtering behavior
type IPFilterConfig struct {
	// Whitelist of allowed IPs/CIDRs (empty = allow all)
	Whitelist []string
	// Blacklist of blocked IPs/CIDRs
	Blacklist []string
	// Error message when blocked
	ErrorMessage string
}

// DefaultIPFilterConfig returns default IP filter config
func DefaultIPFilterConfig() IPFilterConfig {
	return IPFilterConfig{
		Whitelist:     []string{},
		Blacklist:     []string{},
		ErrorMessage:  "Access denied from this IP address",
	}
}

// IPFilterMiddleware creates an IP whitelist/blacklist middleware
func IPFilterMiddleware(config IPFilterConfig) gin.HandlerFunc {
	// Parse CIDR blocks
	whitelistCIDRs := parseCIDRs(config.Whitelist)
	blacklistCIDRs := parseCIDRs(config.Blacklist)

	return func(c *gin.Context) {
		clientIP := GetClientIP(c)
		ip := net.ParseIP(clientIP)
		if ip == nil {
			logger.Component("auth.middleware.ip_filter").
				Warn().
				Str("ip", clientIP).
				Msg("invalid IP address")
			c.Error(errors.NewForbiddenError())
			c.Abort()
			return
		}

		// Check blacklist first
		if isIPBlocked(ip, config.Blacklist, blacklistCIDRs) {
			logger.Component("auth.middleware.ip_filter").
				Warn().
				Str("ip", clientIP).
				Str("path", c.Request.URL.Path).
				Msg("blocked IP attempted access")
			c.Error(errors.NewForbiddenError())
			c.Abort()
			return
		}

		// Check whitelist (if configured)
		if len(config.Whitelist) > 0 {
			if !isIPAllowed(ip, config.Whitelist, whitelistCIDRs) {
				logger.Component("auth.middleware.ip_filter").
					Warn().
					Str("ip", clientIP).
					Str("path", c.Request.URL.Path).
					Msg("non-whitelisted IP attempted access")
				c.Error(errors.NewForbiddenError())
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// parseCIDRs parses CIDR blocks from string list
func parseCIDRs(ipList []string) []*net.IPNet {
	var cidrs []*net.IPNet
	for _, ipStr := range ipList {
		if strings.Contains(ipStr, "/") {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err == nil {
				cidrs = append(cidrs, ipNet)
			}
		}
	}
	return cidrs
}

// isIPBlocked checks if IP is in blacklist
func isIPBlocked(ip net.IP, blacklist []string, cidrs []*net.IPNet) bool {
	// Check exact matches
	ipStr := ip.String()
	for _, blocked := range blacklist {
		if blocked == ipStr {
			return true
		}
	}

	// Check CIDR blocks
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

// isIPAllowed checks if IP is in whitelist
func isIPAllowed(ip net.IP, whitelist []string, cidrs []*net.IPNet) bool {
	// Check exact matches
	ipStr := ip.String()
	for _, allowed := range whitelist {
		if allowed == ipStr {
			return true
		}
	}

	// Check CIDR blocks
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

