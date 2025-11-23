package bootstrap

import (
	"context"
	"fmt"
	"time"

	appcommand "golang-social-media/apps/chat-service/internal/application/command"
	commandcontracts "golang-social-media/apps/chat-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/chat-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/chat-service/internal/application/event_handler"
	chatcache "golang-social-media/apps/chat-service/internal/infrastructure/cache"
	eventbuspublisher "golang-social-media/apps/chat-service/internal/infrastructure/eventbus/publisher"
	eventbussubscriber "golang-social-media/apps/chat-service/internal/infrastructure/eventbus/subscriber"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	domainfactories "golang-social-media/apps/chat-service/internal/domain/factories"
	"golang-social-media/pkg/cache"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	DB                *gorm.DB
	Publisher         *eventbuspublisher.KafkaPublisher
	Cache             cache.Cache
	MessageRepo       *persistence.MessageRepository
	UserRepo          *persistence.UserRepository
	EventDispatcher   *event_dispatcher.Dispatcher
	CreateMessageCmd  commandcontracts.CreateMessageCommand
	UserSubscriber    *eventbussubscriber.UserCreatedSubscriber
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup database
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	// Setup Redis cache
	redisCache, err := setupCache()
	if err != nil {
		logger.Component("chat.bootstrap").
			Warn().
			Err(err).
			Msg("failed to setup cache, continuing without cache")
		redisCache = nil
	}

	// Setup cache wrappers
	var userCache *chatcache.UserCache
	if redisCache != nil {
		userCache = chatcache.NewUserCache(redisCache)
	}

	// Setup mappers
	messageMapper := persistence.NewMessageMapper()

	// Setup repositories
	messageRepo := persistence.NewMessageRepository(db, messageMapper)
	userRepo := persistence.NewUserRepository(db, userCache)

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
	handleUserCreatedCmd := setupHandleUserCreatedCommand(userRepo)

	// Setup subscribers
	userSubscriber, err := setupUserSubscriber(handleUserCreatedCmd)
	if err != nil {
		return nil, err
	}

	logger.Component("chat.bootstrap").
		Info().
		Msg("chat service dependencies initialized")

	return &Dependencies{
		DB:               db,
		Publisher:        publisher,
		Cache:            redisCache,
		MessageRepo:      messageRepo,
		UserRepo:         userRepo,
		EventDispatcher:  eventDispatcher,
		CreateMessageCmd: createMessageCmd,
		UserSubscriber:   userSubscriber,
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

	// Configure connection pool for high concurrency
	sqlDB, err := db.DB()
	if err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to get underlying sql.DB")
		return nil, err
	}

	// Set connection pool settings
	// MaxOpenConns: maximum number of open connections to the database
	// MaxIdleConns: maximum number of connections in the idle connection pool
	// ConnMaxLifetime: maximum amount of time a connection may be reused
	maxOpenConns := config.GetEnvInt("DB_MAX_OPEN_CONNS", 100)
	maxIdleConns := config.GetEnvInt("DB_MAX_IDLE_CONNS", 25)
	connMaxLifetime := config.GetEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 5)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

	logger.Component("chat.bootstrap").
		Info().
		Int("max_open_conns", maxOpenConns).
		Int("max_idle_conns", maxIdleConns).
		Int("conn_max_lifetime_minutes", connMaxLifetime).
		Msg("database connected with connection pool configured")

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

func setupHandleUserCreatedCommand(userRepo *persistence.UserRepository) commandcontracts.HandleUserCreatedCommand {
	handleUserCreatedCmd := appcommand.NewHandleUserCreatedCommand(userRepo)

	logger.Component("chat.bootstrap").
		Info().
		Str("command", "HandleUserCreatedCommand").
		Msg("registered command")

	return handleUserCreatedCmd
}

func setupUserSubscriber(handleUserCreatedCmd commandcontracts.HandleUserCreatedCommand) (*eventbussubscriber.UserCreatedSubscriber, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	groupID := config.GetEnv("CHAT_USER_GROUP_ID", "chat-service-user")

	// Type assertion to get the handler
	handler, ok := handleUserCreatedCmd.(*appcommand.HandleUserCreatedCommandHandler)
	if !ok {
		return nil, fmt.Errorf("invalid HandleUserCreatedCommand type")
	}

	subscriber, err := eventbussubscriber.NewUserCreatedSubscriber(brokers, groupID, handler)
	if err != nil {
		logger.Component("chat.bootstrap").
			Error().
			Err(err).
			Msg("failed to create user subscriber")
		return nil, err
	}

	logger.Component("chat.bootstrap").
		Info().
		Str("subscriber", "UserCreatedSubscriber").
		Str("topic", "user.created").
		Str("group_id", groupID).
		Msg("registered subscriber")

	logger.Component("chat.bootstrap").
		Info().
		Int("total_subscribers", 1).
		Msg("subscribers configured")

	return subscriber, nil
}

func setupCache() (cache.Cache, error) {
	addr := config.GetEnv("REDIS_ADDR", "localhost:6379")
	password := config.GetEnv("REDIS_PASSWORD", "")
	db := config.GetEnvInt("REDIS_DB", 0)

	redisCache, err := cache.NewRedisCache(addr, password, db, "chat.cache")
	if err != nil {
		return nil, err
	}

	logger.Component("chat.bootstrap").
		Info().
		Str("addr", addr).
		Int("db", db).
		Msg("redis cache initialized")

	return redisCache, nil
}

