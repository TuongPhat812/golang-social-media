package bootstrap

import (
	"context"

	appcommand "golang-social-media/apps/chat-service/internal/application/command"
	commandcontracts "golang-social-media/apps/chat-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/chat-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/chat-service/internal/application/event_handler"
	eventbuspublisher "golang-social-media/apps/chat-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	domainfactories "golang-social-media/apps/chat-service/internal/domain/factories"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	DB                *gorm.DB
	Publisher         *eventbuspublisher.KafkaPublisher
	MessageRepo       *persistence.MessageRepository
	EventDispatcher   *event_dispatcher.Dispatcher
	CreateMessageCmd  commandcontracts.CreateMessageCommand
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup database
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	// Setup mappers
	messageMapper := persistence.NewMessageMapper()

	// Setup repositories
	messageRepo := persistence.NewMessageRepository(db, messageMapper)

	// Setup event bus publisher
	publisher, err := setupPublisher()
	if err != nil {
		return nil, err
	}

	// Setup event dispatcher
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup factories
	messageFactory := domainfactories.NewMessageFactory()

	// Setup commands
	createMessageCmd := setupCommands(messageRepo, messageFactory, eventDispatcher)

	logger.Component("chat.bootstrap").
		Info().
		Msg("chat service dependencies initialized")

	return &Dependencies{
		DB:               db,
		Publisher:        publisher,
		MessageRepo:      messageRepo,
		EventDispatcher:  eventDispatcher,
		CreateMessageCmd: createMessageCmd,
	}, nil
}

func setupDatabase() (*gorm.DB, error) {
	dsn := config.GetEnv("CHAT_DATABASE_DSN", "postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect database")
		return nil, err
	}

	logger.Component("chat.bootstrap").
		Info().
		Msg("database connected")

	return db, nil
}

func setupPublisher() (*eventbuspublisher.KafkaPublisher, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	publisher, err := eventbuspublisher.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		return nil, err
	}

	logger.Component("chat.bootstrap").
		Info().
		Msg("kafka publisher initialized")

	return publisher, nil
}

func setupEventDispatcher(publisher *eventbuspublisher.KafkaPublisher) *event_dispatcher.Dispatcher {
	dispatcher := event_dispatcher.NewDispatcher()

	// Create event broker adapter (abstraction over infrastructure)
	eventBrokerAdapter := eventbuspublisher.NewEventBrokerAdapter(publisher)

	// Register MessageCreated handler
	messageCreatedHandler := event_handler.NewMessageCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("MessageCreated", messageCreatedHandler)
	logger.Component("chat.bootstrap").
		Info().
		Str("event_type", "MessageCreated").
		Str("handler", "MessageCreatedHandler").
		Msg("registered event handler")

	logger.Component("chat.bootstrap").
		Info().
		Int("total_handlers", 1).
		Msg("event dispatcher configured")

	return dispatcher
}

func setupCommands(
	messageRepo *persistence.MessageRepository,
	messageFactory domainfactories.MessageFactory,
	eventDispatcher *event_dispatcher.Dispatcher,
) commandcontracts.CreateMessageCommand {
	createMessageCmd := appcommand.NewCreateMessageCommand(messageRepo, messageFactory, eventDispatcher)

	logger.Component("chat.bootstrap").
		Info().
		Str("command", "CreateMessageCommand").
		Msg("registered command")

	logger.Component("chat.bootstrap").
		Info().
		Int("total_commands", 1).
		Msg("commands configured")

	return createMessageCmd
}

