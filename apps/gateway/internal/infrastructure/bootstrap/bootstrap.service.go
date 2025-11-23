package bootstrap

import (
	"context"
	"io"
	"strings"

	appcommand "golang-social-media/apps/gateway/internal/application/command"
	commandcontracts "golang-social-media/apps/gateway/internal/application/command/contracts"
	appquery "golang-social-media/apps/gateway/internal/application/query"
	querycontracts "golang-social-media/apps/gateway/internal/application/query/contracts"
	authclient "golang-social-media/apps/gateway/internal/infrastructure/auth"
	authgrpc "golang-social-media/apps/gateway/internal/infrastructure/grpc/auth"
	chatclient "golang-social-media/apps/gateway/internal/infrastructure/grpc/chat"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	ChatClient         *chatclient.Client
	AuthClient         *authclient.Client
	AuthGRPCClient     *authgrpc.Client
	CreateMessageCmd   commandcontracts.CreateMessageCommand
	RegisterUserCmd    commandcontracts.RegisterUserCommand
	LoginUserCmd       commandcontracts.LoginUserCommand
	GetUserProfileQuery querycontracts.GetUserProfileQuery
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup Gin mode
	setupGinMode()

	// Setup clients
	chatClient, err := setupChatClient(ctx)
	if err != nil {
		return nil, err
	}

	authClient := setupAuthClient()

	// Setup auth gRPC client for token validation
	authGRPCClient, err := setupAuthGRPCClient(ctx)
	if err != nil {
		return nil, err
	}

	// Setup commands
	createMessageCmd := appcommand.NewCreateMessageCommand(chatClient)
	registerUserCmd := appcommand.NewRegisterUserCommand(authClient)
	loginUserCmd := appcommand.NewLoginUserCommand(authClient)

	// Setup queries
	getUserProfileQuery := appquery.NewGetUserProfileQuery(authClient)

	logger.Component("gateway.bootstrap").
		Info().
		Msg("gateway service dependencies initialized")

	return &Dependencies{
		ChatClient:         chatClient,
		AuthClient:         authClient,
		AuthGRPCClient:     authGRPCClient,
		CreateMessageCmd:   createMessageCmd,
		RegisterUserCmd:    registerUserCmd,
		LoginUserCmd:       loginUserCmd,
		GetUserProfileQuery: getUserProfileQuery,
	}, nil
}

func setupGinMode() {
	switch strings.ToLower(config.GetEnv("GIN_MODE", "debug")) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if strings.EqualFold(config.GetEnv("GIN_DISABLE_ACCESS_LOG", "false"), "true") {
		gin.DefaultWriter = io.Discard
	}

	logger.Component("gateway.bootstrap").
		Info().
		Str("mode", gin.Mode()).
		Bool("access_log_disabled", gin.DefaultWriter == io.Discard).
		Msg("Gin configured")
}

func setupChatClient(ctx context.Context) (*chatclient.Client, error) {
	client, err := chatclient.NewClient(ctx)
	if err != nil {
		logger.Component("gateway.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect to chat service")
		return nil, err
	}

	logger.Component("gateway.bootstrap").
		Info().
		Msg("chat service client connected")

	return client, nil
}

func setupAuthClient() *authclient.Client {
	authBaseURL := config.GetEnv("AUTH_SERVICE_URL", "http://localhost:9101")
	client := authclient.NewClient(authBaseURL)

	logger.Component("gateway.bootstrap").
		Info().
		Str("base_url", authBaseURL).
		Msg("auth service HTTP client initialized")

	return client
}

func setupAuthGRPCClient(ctx context.Context) (*authgrpc.Client, error) {
	client, err := authgrpc.NewClient(ctx)
	if err != nil {
		logger.Component("gateway.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect to auth service gRPC")
		return nil, err
	}

	logger.Component("gateway.bootstrap").
		Info().
		Msg("auth service gRPC client connected")

	return client, nil
}

