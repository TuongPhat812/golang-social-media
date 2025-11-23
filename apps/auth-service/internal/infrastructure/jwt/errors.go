package jwt

import "errors"

var (
	ErrTokenBlacklisted = errors.New("token is blacklisted")
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token expired")
)

