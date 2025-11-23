package jwt

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetTokenID generates a unique ID for a token (for blacklisting)
// Uses SHA256 hash of the token string
func GetTokenID(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}

