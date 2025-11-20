package bootstrap

import (
	"context"

	appcommand "golang-social-media/apps/auth-service/internal/application/command"
	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	appquery "golang-social-media/apps/auth-service/internal/application/query"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/auth-service/internal/application/event_handler"
	eventbuspublisher "golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	domainfactories "golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	Publisher         *eventbuspublisher.KafkaPublisher
	UserRepo          *memory.UserRepository
	EventDispatcher   *event_dispatcher.Dispatcher
	RegisterUserCmd   commandcontracts.RegisterUserCommand
	LoginUserCmd      *appcommand.LoginUserHandler
	GetUserProfileQuery querycontracts.GetUserProfileQuery
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup repositories
	userRepo := memory.NewUserRepository()

	// Setup event bus publisher
	publisher, err := setupPublisher()
	if err != nil {
		return nil, err
	}

	// Setup event dispatcher
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup factories
	userFactory := domainfactories.NewUserFactory()

	// Setup commands
	registerUserCmd := appcommand.NewRegisterUserCommand(userRepo, userFactory, eventDispatcher)
	tokenStore := appcommand.NewTokenStore()
	loginUserCmd := appcommand.NewLoginUserHandler(userRepo, tokenStore)

	// Setup queries
	getUserProfileQuery := appquery.NewGetUserProfileHandler(userRepo)

	logger.Component("auth.bootstrap").
		Info().
		Msg("auth service dependencies initialized")

	return &Dependencies{
		Publisher:         publisher,
		UserRepo:         userRepo,
		EventDispatcher:  eventDispatcher,
		RegisterUserCmd:   registerUserCmd,
		LoginUserCmd:      loginUserCmd,
		GetUserProfileQuery: getUserProfileQuery,
	}, nil
}

func setupPublisher() (*eventbuspublisher.KafkaPublisher, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	publisher, err := eventbuspublisher.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("auth.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		return nil, err
	}

	logger.Component("auth.bootstrap").
		Info().
		Msg("kafka publisher initialized")

	return publisher, nil
}

func setupEventDispatcher(publisher *eventbuspublisher.KafkaPublisher) *event_dispatcher.Dispatcher {
	dispatcher := event_dispatcher.NewDispatcher()

	// Create event broker adapter (abstraction over infrastructure)
	eventBrokerAdapter := eventbuspublisher.NewEventBrokerAdapter(publisher)

	// Register UserCreated handler
	userCreatedHandler := event_handler.NewUserCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("UserCreated", userCreatedHandler)
	logger.Component("auth.bootstrap").
		Info().
		Str("event_type", "UserCreated").
		Str("handler", "UserCreatedHandler").
		Msg("registered event handler")

	logger.Component("auth.bootstrap").
		Info().
		Int("total_handlers", 1).
		Msg("event dispatcher configured")

	return dispatcher
}

