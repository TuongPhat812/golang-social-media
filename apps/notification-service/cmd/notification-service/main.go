package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang-social-media/apps/notification-service/internal/application/notifications"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus"
	grpcserver "golang-social-media/apps/notification-service/internal/infrastructure/grpc"
	interfaces "golang-social-media/apps/notification-service/internal/interfaces/grpc"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})

	publisher, err := eventbus.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create notification kafka publisher")
		os.Exit(1)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close notification kafka publisher")
		}
	}()

	notificationService := notifications.NewService(publisher)

	subscriber, err := eventbus.NewSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_CHAT_GROUP_ID", "notification-service-chat"),
		notificationService,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create notification kafka subscriber")
		os.Exit(1)
	}
	defer func() {
		if err := subscriber.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close notification kafka subscriber")
		}
	}()

	subscriber.ConsumeChatCreated(ctx)

	logger.Info().Msg("notification service ready")

	port := config.GetEnvInt("NOTIFICATION_SERVICE_PORT", 9100)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		interfaces.Register(server, notificationService)
	}); err != nil {
		logger.Error().Err(err).Msg("failed to serve notification gRPC")
		os.Exit(1)
	}
}
