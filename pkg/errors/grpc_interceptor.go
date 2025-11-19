package errors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"golang-social-media/pkg/logger"
	"github.com/rs/zerolog"
)

// GRPCErrorInterceptor creates a gRPC unary interceptor that transforms errors
func GRPCErrorInterceptor(transformer *Transformer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Transform error
		appErr := transformer.Transform(err)

		// Log the error
		log := logger.Component("grpc.error_interceptor")
		logGRPCError(log, appErr, info)

		// Convert to gRPC status
		grpcCode := mapHTTPStatusToGRPCCode(appErr.HTTPStatus)
		// Include error code in message for client to parse
		message := appErr.Message
		if appErr.Code != "" {
			message = string(appErr.Code) + ": " + message
		}
		grpcStatus := status.New(grpcCode, message)

		return nil, grpcStatus.Err()
	}
}

func logGRPCError(log *zerolog.Logger, appErr *AppError, info *grpc.UnaryServerInfo) {
	event := log.Error().
		Str("error_code", string(appErr.Code)).
		Str("error_message", appErr.Message).
		Str("method", info.FullMethod)

	if appErr.Original != nil {
		event = event.Err(appErr.Original)
	}

	if len(appErr.Details) > 0 {
		event = event.Interface("details", appErr.Details)
	}

	event.Msg("gRPC error")
}

// mapHTTPStatusToGRPCCode maps HTTP status codes to gRPC status codes
func mapHTTPStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case 200:
		return codes.OK
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 408:
		return codes.DeadlineExceeded
	case 409:
		return codes.AlreadyExists
	case 429:
		return codes.ResourceExhausted
	case 500:
		return codes.Internal
	case 503:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

