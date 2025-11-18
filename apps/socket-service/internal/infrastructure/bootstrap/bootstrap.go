package bootstrap

import (
	"context"

	appevents "golang-social-media/apps/socket-service/internal/application/events"
	eventbussubscriber "golang-social-media/apps/socket-service/internal/infrastructure/eventbus/subscriber"
	"golang-social-media/apps/socket-service/internal/interfaces/socket"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	Hub                      *socket.Hub
	EventService             appevents.Service
	ChatSubscriber           *eventbussubscriber.ChatCreatedSubscriber
	NotificationSubscriber   *eventbussubscriber.NotificationCreatedSubscriber
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup socket hub
	hub := socket.NewHub()

	// Setup event service
	eventService := appevents.NewService(hub)

	// Setup subscribers
	chatSubscriber, err := setupChatSubscriber(eventService)
	if err != nil {
		return nil, err
	}

	notificationSubscriber, err := setupNotificationSubscriber(eventService)
	if err != nil {
		return nil, err
	}

	logger.Component("socket.bootstrap").
		Info().
		Msg("socket service dependencies initialized")

	return &Dependencies{
		Hub:                    hub,
		EventService:           eventService,
		ChatSubscriber:         chatSubscriber,
		NotificationSubscriber: notificationSubscriber,
	}, nil
}

func setupChatSubscriber(eventService appevents.Service) (*eventbussubscriber.ChatCreatedSubscriber, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	groupID := config.GetEnv("SOCKET_CHAT_GROUP_ID", "socket-service-chat")

	subscriber, err := eventbussubscriber.NewChatCreatedSubscriber(brokers, groupID, eventService)
	if err != nil {
		logger.Component("socket.bootstrap").
			Error().
			Err(err).
			Msg("failed to create chat subscriber")
		return nil, err
	}

	logger.Component("socket.bootstrap").
		Info().
		Str("subscriber", "ChatCreatedSubscriber").
		Str("topic", "chat.created").
		Msg("registered subscriber")

	return subscriber, nil
}

func setupNotificationSubscriber(eventService appevents.Service) (*eventbussubscriber.NotificationCreatedSubscriber, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	groupID := config.GetEnv("SOCKET_NOTIFICATION_GROUP_ID", "socket-service-notification")

	subscriber, err := eventbussubscriber.NewNotificationCreatedSubscriber(brokers, groupID, eventService)
	if err != nil {
		logger.Component("socket.bootstrap").
			Error().
			Err(err).
			Msg("failed to create notification subscriber")
		return nil, err
	}

	logger.Component("socket.bootstrap").
		Info().
		Str("subscriber", "NotificationCreatedSubscriber").
		Str("topic", "notification.created").
		Msg("registered subscriber")

	return subscriber, nil
}
