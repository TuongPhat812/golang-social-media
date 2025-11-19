package random

import (
	"crypto/rand"
	"encoding/hex"
)

func String(n int) string {
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "token-fallback"
	}
	return hex.EncodeToString(b)[:n*2]
}
