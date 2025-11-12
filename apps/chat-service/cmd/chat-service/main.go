package main

import (
	"fmt"
	"os"

	"golang-social-media/apps/chat-service/internal/application/messages"
	"golang-social-media/apps/chat-service/internal/infrastructure/eventbus"
	grpcserver "golang-social-media/apps/chat-service/internal/infrastructure/grpc"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	chatgrpc "golang-social-media/apps/chat-service/internal/interfaces/grpc/chat"
	"golang-social-media/pkg/config"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"golang-social-media/pkg/logger"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()

	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	publisher, err := eventbus.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create kafka publisher")
		os.Exit(1)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close kafka publisher")
		}
	}()

	dsn := config.GetEnv("CHAT_DATABASE_DSN", "postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect database")
		os.Exit(1)
	}
	messageRepository := persistence.NewMessageRepository(db)
	messageService := messages.NewService(messageRepository, publisher)

	port := config.GetEnvInt("CHAT_SERVICE_PORT", 9000)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		handler := chatgrpc.NewHandler(messageService)
		chatv1.RegisterChatServiceServer(server, handler)
	}); err != nil {
		logger.Error().Err(err).Msg("failed to start gRPC server")
		os.Exit(1)
	}
}
