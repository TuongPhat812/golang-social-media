package main

import (
	"fmt"
	"os"

	"golang-social-media/apps/auth-service/internal/application/command"
	"golang-social-media/apps/auth-service/internal/application/query"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/interfaces/rest"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
)

func main() {
	logger.SetModule("auth-service")
	config.LoadEnv()

	repo := memory.NewUserRepository()
	tokenStore := command.NewTokenStore()

	registerHandler := command.NewRegisterUserHandler(repo, nil)
	loginHandler := command.NewLoginUserHandler(repo, tokenStore)
	getProfileHandler := query.NewGetUserProfileHandler(repo)

	router := rest.NewRouter(rest.Handlers{
		RegisterUser: registerHandler,
		LoginUser:    loginHandler,
		GetProfile:   getProfileHandler,
	})

	port := config.GetEnvInt("AUTH_SERVICE_PORT", 9101)
	addr := fmt.Sprintf(":%d", port)

	logger.Component("auth.http").
		Info().
		Str("addr", addr).
		Msg("auth service starting")

	if err := router.Run(addr); err != nil {
		logger.Component("auth.http").
			Error().
			Err(err).
			Msg("auth service failed")
		os.Exit(1)
	}
}
