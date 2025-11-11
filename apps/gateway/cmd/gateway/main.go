package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/myself/golang-social-media/pkg/config"
	"github.com/myself/golang-social-media/apps/gateway/internal/application/messages"
	"github.com/myself/golang-social-media/apps/gateway/internal/application/users"
	chatclient "github.com/myself/golang-social-media/apps/gateway/internal/infrastructure/grpc/chat"
	httpserver "github.com/myself/golang-social-media/apps/gateway/internal/infrastructure/http"
)

func main() {
	config.LoadEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	chatClient, err := chatclient.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to connect to chat service: %v", err)
	}
	defer func() {
		if err := chatClient.Close(); err != nil {
			log.Printf("failed to close chat client: %v", err)
		}
	}()

	userService := users.NewService()
	messageService := messages.NewService(chatClient)

	router := httpserver.NewRouter(userService, messageService)

	port := config.GetEnvInt("GATEWAY_PORT", 8080)
	addr := fmt.Sprintf(":%d", port)

	log.Printf("gateway service starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}
}
