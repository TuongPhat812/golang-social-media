package contracts

import "context"

type RevokeTokenCommandRequest struct {
	Token string // Token to revoke (can be access or refresh token)
}

type RevokeTokenCommand interface {
	Execute(ctx context.Context, req RevokeTokenCommandRequest) error
}

