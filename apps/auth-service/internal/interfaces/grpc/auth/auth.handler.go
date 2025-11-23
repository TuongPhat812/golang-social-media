package auth

import (
	"context"
	"time"

	bootstrap "golang-social-media/apps/auth-service/internal/infrastructure/bootstrap"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/pkg/logger"
	authv1 "golang-social-media/pkg/gen/auth/v1"
)

type Handler struct {
	validateTokenQuery querycontracts.ValidateTokenQuery
	authv1.UnimplementedAuthServiceServer
}

func NewHandler(deps *bootstrap.Dependencies) *Handler {
	return &Handler{
		validateTokenQuery: deps.ValidateTokenQuery,
	}
}

func (h *Handler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	startTime := time.Now()

	token := req.GetToken()
	if token == "" {
		logger.Component("auth.grpc.validate_token").
			Warn().
			Msg("empty token in request")
		return &authv1.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
		}, nil
	}

	// Execute query
	queryStart := time.Now()
	resp, err := h.validateTokenQuery.Execute(ctx, token)
	queryDuration := time.Since(queryStart)

	if err != nil {
		totalDuration := time.Since(startTime)
		logger.Component("auth.grpc.validate_token").
			Error().
			Err(err).
			Dur("query_exec_ms", queryDuration).
			Dur("total_ms", totalDuration).
			Msg("failed to validate token")
		return &authv1.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
		}, nil // Return valid=false, not error
	}

	totalDuration := time.Since(startTime)

	logger.Component("auth.grpc.validate_token").
		Info().
		Bool("valid", resp.Valid).
		Str("user_id", resp.UserID).
		Dur("query_exec_ms", queryDuration).
		Dur("total_ms", totalDuration).
		Msg("gRPC request completed")

	return &authv1.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserId: resp.UserID,
	}, nil
}

