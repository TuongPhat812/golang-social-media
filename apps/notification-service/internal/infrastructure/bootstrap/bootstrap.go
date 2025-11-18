package bootstrap

import (
	"context"

	command "golang-social-media/apps/notification-service/internal/application/command"
	commandcontracts "golang-social-media/apps/notification-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/notification-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/notification-service/internal/application/event_handler"
	query "golang-social-media/apps/notification-service/internal/application/query"
	querycontracts "golang-social-media/apps/notification-service/internal/application/query/contracts"
	eventbuspublisher "golang-social-media/apps/notification-service/internal/infrastructure/eventbus/publisher"
	eventbussubscriber "golang-social-media/apps/notification-service/internal/infrastructure/eventbus/subscriber"
	scylladb "golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"

	"github.com/gocql/gocql"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	Publisher               *eventbuspublisher.KafkaPublisher
	Session                 *gocql.Session
	NotificationRepo        *scylladb.NotificationRepository
	UserRepo                *scylladb.UserRepository
	EventDispatcher         *event_dispatcher.Dispatcher
	CreateNotificationCmd   commandcontracts.CreateNotificationCommand
	MarkNotificationReadCmd commandcontracts.MarkNotificationReadCommand
	GetNotificationsQuery   querycontracts.GetNotificationsQuery
	HandleChatCreatedCmd    *command.HandleChatCreatedCommandHandler
	HandleUserCreatedCmd    *command.HandleUserCreatedCommandHandler
	ChatSubscriber          *eventbussubscriber.ChatCreatedSubscriber
	UserSubscriber          *eventbussubscriber.UserCreatedSubscriber
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})

	// Setup Kafka Publisher
	publisher, err := eventbuspublisher.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		return nil, err
	}

	// Setup ScyllaDB
	scyllaHosts := config.GetEnvStringSlice("SCYLLA_HOSTS", []string{"localhost:9042"})
	scyllaKeyspace := config.GetEnv("SCYLLA_KEYSPACE", "notification_service")
	session, err := scylladb.NewSession(scyllaHosts, scyllaKeyspace)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Strs("hosts", scyllaHosts).
			Str("keyspace", scyllaKeyspace).
			Msg("failed to connect scylla")
		return nil, err
	}

	notificationRepo := scylladb.NewNotificationRepository(session)
	userRepo := scylladb.NewUserRepository(session)

	// Setup event dispatcher and handlers
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup commands
	commands := setupCommands(notificationRepo, userRepo, eventDispatcher)

	// Setup queries
	queries := setupQueries(notificationRepo)

	// Setup subscribers
	subscribers, err := setupSubscribers(ctx, brokers, commands)
	if err != nil {
		session.Close()
		return nil, err
	}

	return &Dependencies{
		Publisher:               publisher,
		Session:                 session,
		NotificationRepo:        notificationRepo,
		UserRepo:                userRepo,
		EventDispatcher:         eventDispatcher,
		CreateNotificationCmd:   commands.CreateNotification,
		MarkNotificationReadCmd: commands.MarkNotificationRead,
		GetNotificationsQuery:   queries.GetNotifications,
		HandleChatCreatedCmd:    commands.HandleChatCreated,
		HandleUserCreatedCmd:    commands.HandleUserCreated,
		ChatSubscriber:          subscribers.Chat,
		UserSubscriber:          subscribers.User,
	}, nil
}

// setupEventDispatcher configures the event dispatcher with all handlers
func setupEventDispatcher(publisher *eventbuspublisher.KafkaPublisher) *event_dispatcher.Dispatcher {
	dispatcher := event_dispatcher.NewDispatcher()

	// Create event broker adapter (abstraction over infrastructure)
	eventBrokerAdapter := eventbuspublisher.NewEventBrokerAdapter(publisher)

	// Register NotificationCreated handler
	notificationCreatedHandler := event_handler.NewNotificationCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("NotificationCreated", notificationCreatedHandler)
	logger.Component("notification.bootstrap").
		Info().
		Str("event_type", "NotificationCreated").
		Str("handler", "NotificationCreatedHandler").
		Msg("registered event handler")

	// Register NotificationRead handler
	notificationReadHandler := event_handler.NewNotificationReadHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("NotificationRead", notificationReadHandler)
	logger.Component("notification.bootstrap").
		Info().
		Str("event_type", "NotificationRead").
		Str("handler", "NotificationReadHandler").
		Msg("registered event handler")

	logger.Component("notification.bootstrap").
		Info().
		Int("total_handlers", 2).
		Msg("event dispatcher configured")

	return dispatcher
}

type commands struct {
	CreateNotification   commandcontracts.CreateNotificationCommand
	MarkNotificationRead commandcontracts.MarkNotificationReadCommand
	HandleChatCreated    *command.HandleChatCreatedCommandHandler
	HandleUserCreated    *command.HandleUserCreatedCommandHandler
}

// setupCommands initializes all command handlers
func setupCommands(
	notificationRepo *scylladb.NotificationRepository,
	userRepo *scylladb.UserRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
) commands {
	createNotificationCmd := command.NewCreateNotificationCommand(notificationRepo, eventDispatcher)
	markNotificationReadCmd := command.NewMarkNotificationReadCommand(notificationRepo, eventDispatcher)
	handleChatCreatedCmd := command.NewHandleChatCreatedCommand(createNotificationCmd)
	handleUserCreatedCmd := command.NewHandleUserCreatedCommand(userRepo, createNotificationCmd)

	logger.Component("notification.bootstrap").
		Info().
		Str("command", "CreateNotificationCommand").
		Msg("registered command")

	logger.Component("notification.bootstrap").
		Info().
		Str("command", "MarkNotificationReadCommand").
		Msg("registered command")

	logger.Component("notification.bootstrap").
		Info().
		Str("command", "HandleChatCreatedCommand").
		Msg("registered command")

	logger.Component("notification.bootstrap").
		Info().
		Str("command", "HandleUserCreatedCommand").
		Msg("registered command")

	logger.Component("notification.bootstrap").
		Info().
		Int("total_commands", 4).
		Msg("commands configured")

	return commands{
		CreateNotification:   createNotificationCmd,
		MarkNotificationRead: markNotificationReadCmd,
		HandleChatCreated:    handleChatCreatedCmd,
		HandleUserCreated:    handleUserCreatedCmd,
	}
}

type queries struct {
	GetNotifications querycontracts.GetNotificationsQuery
}

// setupQueries initializes all query handlers
func setupQueries(notificationRepo *scylladb.NotificationRepository) queries {
	getNotificationsQuery := query.NewGetNotificationsQuery(notificationRepo)

	logger.Component("notification.bootstrap").
		Info().
		Str("query", "GetNotificationsQuery").
		Msg("registered query")

	logger.Component("notification.bootstrap").
		Info().
		Int("total_queries", 1).
		Msg("queries configured")

	return queries{
		GetNotifications: getNotificationsQuery,
	}
}

type subscribers struct {
	Chat *eventbussubscriber.ChatCreatedSubscriber
	User *eventbussubscriber.UserCreatedSubscriber
}

// setupSubscribers initializes all event subscribers
func setupSubscribers(
	ctx context.Context,
	brokers []string,
	commands commands,
) (subscribers, error) {
	chatSubscriber, err := eventbussubscriber.NewChatCreatedSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_CHAT_GROUP_ID", "notification-service-chat"),
		commands.HandleChatCreated,
	)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create chat subscriber")
		return subscribers{}, err
	}

	userSubscriber, err := eventbussubscriber.NewUserCreatedSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_USER_GROUP_ID", "notification-service-user"),
		commands.HandleUserCreated,
	)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create user subscriber")
		return subscribers{}, err
	}

	logger.Component("notification.bootstrap").
		Info().
		Str("subscriber", "ChatCreatedSubscriber").
		Str("topic", events.TopicChatCreated).
		Msg("registered subscriber")

	logger.Component("notification.bootstrap").
		Info().
		Str("subscriber", "UserCreatedSubscriber").
		Str("topic", events.TopicUserCreated).
		Msg("registered subscriber")

	logger.Component("notification.bootstrap").
		Info().
		Int("total_subscribers", 2).
		Msg("subscribers configured")

	return subscribers{
		Chat: chatSubscriber,
		User: userSubscriber,
	}, nil
}
