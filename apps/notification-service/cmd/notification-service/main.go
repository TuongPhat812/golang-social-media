package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang-social-media/pkg/config"
	"golang-social-media/apps/notification-service/internal/application/notifications"
	"golang-social-media/apps/notification-service/internal/infrastructure/eventbus"
	grpcserver "golang-social-media/apps/notification-service/internal/infrastructure/grpc"
	interfaces "golang-social-media/apps/notification-service/internal/interfaces/grpc"
	"google.golang.org/grpc"
)

func main() {
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})

	publisher, err := eventbus.NewKafkaPublisher(brokers)
	if err != nil {
		log.Fatalf("failed to create kafka publisher: %v", err)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("failed to close kafka publisher: %v", err)
		}
	}()

	notificationService := notifications.NewService(publisher)

	subscriber, err := eventbus.NewSubscriber(
		brokers,
		config.GetEnv("NOTIFICATION_CHAT_GROUP_ID", "notification-service-chat"),
		notificationService,
	)
	if err != nil {
		log.Fatalf("failed to create kafka subscriber: %v", err)
	}
	defer func() {
		if err := subscriber.Close(); err != nil {
			log.Printf("failed to close kafka subscriber: %v", err)
		}
	}()

	subscriber.ConsumeChatCreated(ctx)

	sample := notificationService.SampleUser()
	log.Printf("notification service ready with sample user: %+v", sample)

	port := config.GetEnvInt("NOTIFICATION_SERVICE_PORT", 9100)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		interfaces.Register(server, notificationService)
	}); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
