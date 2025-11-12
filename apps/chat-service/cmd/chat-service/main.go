package main

import (
	"fmt"
	"log"

	"golang-social-media/apps/chat-service/internal/application/messages"
	"golang-social-media/apps/chat-service/internal/infrastructure/eventbus"
	grpcserver "golang-social-media/apps/chat-service/internal/infrastructure/grpc"
	chatgrpc "golang-social-media/apps/chat-service/internal/interfaces/grpc/chat"
	"golang-social-media/pkg/config"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"google.golang.org/grpc"
)

func main() {
	config.LoadEnv()

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

	messageService := messages.NewService(publisher)
	sample := messageService.SampleMessage()
	log.Printf("chat service ready with sample message: %+v", sample)

	port := config.GetEnvInt("CHAT_SERVICE_PORT", 9000)
	addr := fmt.Sprintf(":%d", port)

	if err := grpcserver.Start(addr, func(server *grpc.Server) {
		handler := chatgrpc.NewHandler(messageService)
		chatv1.RegisterChatServiceServer(server, handler)
	}); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
