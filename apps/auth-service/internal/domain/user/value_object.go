package user

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// TokenID represents a token identifier value object
type TokenID struct {
	value string
}

// NewTokenID creates a new TokenID from a token string
// Uses SHA256 hash of the token as the ID
func NewTokenID(token string) TokenID {
	hash := sha256.Sum256([]byte(token))
	return TokenID{
		value: hex.EncodeToString(hash[:]),
	}
}

// String returns the string representation of TokenID
func (t TokenID) String() string {
	return t.value
}

// Password represents a password value object
type Password struct {
	hashed string
}

// NewPassword creates a new Password value object
// In production, should use bcrypt or similar
func NewPassword(plainPassword string) Password {
	// Simple hash for now - should use bcrypt in production
	hash := sha256.Sum256([]byte(plainPassword))
	return Password{
		hashed: hex.EncodeToString(hash[:]),
	}
}

// String returns the hashed password
func (p Password) String() string {
	return p.hashed
}

// Verify checks if a plain password matches the hashed password
func (p Password) Verify(plainPassword string) bool {
	hash := sha256.Sum256([]byte(plainPassword))
	hashed := hex.EncodeToString(hash[:])
	return p.hashed == hashed
}

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object
func NewEmail(email string) Email {
	return Email{value: email}
}

// String returns the email string
func (e Email) String() string {
	return e.value
}

// Name represents a user name value object
type Name struct {
	value string
}

// NewName creates a new Name value object
func NewName(name string) Name {
	return Name{value: name}
}

// String returns the name string
func (n Name) String() string {
	return n.value
}

// Timestamp represents a timestamp value object
type Timestamp struct {
	value time.Time
}

// NewTimestamp creates a new Timestamp value object
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Time returns the time value
func (t Timestamp) Time() time.Time {
	return t.value
}

// String returns the RFC3339 string representation
func (t Timestamp) String() string {
	return t.value.UTC().Format(time.RFC3339)
}

