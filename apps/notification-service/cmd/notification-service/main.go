package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	command "golang-social-media/apps/notification-service/internal/application/command"
	"golang-social-media/apps/notification-service/internal/application/consumers"
	"golang-social-media/apps/notification-service/internal/application/notifications"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus"
	grpcserver "golang-social-media/apps/notification-service/internal/infrastructure/grpc"
	scylladb "golang-social-media/apps/notification-service/internal/infrastructure/persistence/scylla"
	interfaces "golang-social-media/apps/notification-service/internal/interfaces/grpc"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger.SetModule("notification-service")
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})

	publisher, err := eventbus.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		os.Exit(1)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close kafka publisher")
		}
	}()

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
		os.Exit(1)
	}
	defer session.Close()

	notificationRepo := scylladb.NewNotificationRepository(session)
	userRepo := scylladb.NewUserRepository(session)

	createNotificationCmd := command.NewCreateNotificationCommand(notificationRepo, publisher)
	notificationService := notifications.NewService(createNotificationCmd)

	userConsumer := consumers.NewUserCreatedConsumer(userRepo, createNotificationCmd)

	chatSubscriber, err := eventbus.NewSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_CHAT_GROUP_ID", "notification-service-chat"),
		notificationService,
	)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create chat subscriber")
		os.Exit(1)
	}
	defer func() {
		if err := chatSubscriber.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close chat subscriber")
		}
	}()
	chatSubscriber.ConsumeChatCreated(ctx)

	userSubscriber, err := eventbus.NewUserSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_USER_GROUP_ID", "notification-service-user"),
		userConsumer,
	)
	if err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to create user subscriber")
		os.Exit(1)
	}
	defer func() {
		if err := userSubscriber.Close(); err != nil {
			logger.Component("notification.bootstrap").
				Error().
				Err(err).
				Msg("failed to close user subscriber")
		}
	}()
	userSubscriber.Consume(ctx)

	logger.Component("notification.bootstrap").
		Info().
		Msg("notification service ready")

	port := config.GetEnvInt("NOTIFICATION_SERVICE_PORT", 9100)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		interfaces.Register(server, notificationService)
	}); err != nil {
		logger.Component("notification.bootstrap").
			Error().
			Err(err).
			Msg("failed to serve notification gRPC")
		os.Exit(1)
	}
}
