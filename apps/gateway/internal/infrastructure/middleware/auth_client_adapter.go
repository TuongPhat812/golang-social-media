package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	authgrpc "golang-social-media/apps/gateway/internal/infrastructure/grpc/auth"
	authv1 "golang-social-media/pkg/gen/auth/v1"
)

// AuthGRPCClientAdapter adapts auth gRPC client to AuthClient interface
type AuthGRPCClientAdapter struct {
	client *authgrpc.Client
}

func NewAuthGRPCClientAdapter(client *authgrpc.Client) *AuthGRPCClientAdapter {
	return &AuthGRPCClientAdapter{
		client: client,
	}
}

func (a *AuthGRPCClientAdapter) ValidateToken(c gin.Context, token string) (userID string, valid bool, err error) {
	ctx := c.Request.Context()
	resp, err := a.client.ValidateToken(ctx, token)
	if err != nil {
		return "", false, err
	}
	return resp.GetUserId(), resp.GetValid(), nil
}

