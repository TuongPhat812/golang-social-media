package bootstrap

import (
	"context"
	"time"

	appcommand "golang-social-media/apps/auth-service/internal/application/command"
	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	appquery "golang-social-media/apps/auth-service/internal/application/query"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/auth-service/internal/application/event_handler"
	authcache "golang-social-media/apps/auth-service/internal/infrastructure/cache"
	eventbuspublisher "golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
	redispersistence "golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	domainfactories "golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/pkg/cache"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	DB                  *gorm.DB
	Publisher           *eventbuspublisher.KafkaPublisher
	Cache               cache.Cache
	UserRepo            *postgres.UserRepository
	RoleRepo            *postgres.RoleRepository
	PermissionRepo      *postgres.PermissionRepository
	UserRoleRepo        *postgres.UserRoleRepository
	RolePermissionRepo  *postgres.RolePermissionRepository
	TokenBlacklistRepo  *redispersistence.TokenBlacklistRepository
	EventDispatcher     *event_dispatcher.Dispatcher
	JwtService          *jwt.Service
	RegisterUserCmd     commandcontracts.RegisterUserCommand
	LoginUserCmd        *appcommand.LoginUserHandler
	LogoutUserCmd       commandcontracts.LogoutUserCommand
	RefreshTokenCmd     commandcontracts.RefreshTokenCommand
	RevokeTokenCmd      commandcontracts.RevokeTokenCommand
	UpdateProfileCmd    commandcontracts.UpdateProfileCommand
	ChangePasswordCmd   commandcontracts.ChangePasswordCommand
	GetUserProfileQuery querycontracts.GetUserProfileQuery
	GetCurrentUserQuery querycontracts.GetCurrentUserQuery
	ValidateTokenQuery  querycontracts.ValidateTokenQuery
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup PostgreSQL database
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	// Setup Redis cache
	redisCache, err := setupCache()
	if err != nil {
		logger.Component("auth.bootstrap").
			Warn().
			Err(err).
			Msg("failed to setup cache, continuing without cache")
		redisCache = nil
	}

	// Setup cache wrappers
	var userCache *authcache.UserCache
	if redisCache != nil {
		userCache = authcache.NewUserCache(redisCache)
	}

	// Setup repositories
	userRepo := postgres.NewUserRepository(db, userCache)
	roleRepo := postgres.NewRoleRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	userRoleRepo := postgres.NewUserRoleRepository(db)
	rolePermissionRepo := postgres.NewRolePermissionRepository(db)

	// Setup token blacklist repository
	var tokenBlacklistRepo *redispersistence.TokenBlacklistRepository
	if redisCache != nil {
		tokenBlacklistRepo = redispersistence.NewTokenBlacklistRepository(redisCache)
	}

	// Setup event bus publisher
	publisher, err := setupPublisher()
	if err != nil {
		return nil, err
	}

	// Setup event dispatcher
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup factories
	userFactory := domainfactories.NewUserFactory()

	// Setup JWT service
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key-change-in-production")
	accessExpirationHours := config.GetEnvInt("JWT_ACCESS_EXPIRATION_HOURS", 1)   // Default 1 hour
	refreshExpirationHours := config.GetEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 168) // Default 7 days
	jwtService := jwt.NewService(jwtSecret, accessExpirationHours, refreshExpirationHours)

	// Setup commands
	registerUserCmd := appcommand.NewRegisterUserCommand(userRepo, userFactory, eventDispatcher)
	loginUserCmd := appcommand.NewLoginUserHandler(userRepo, jwtService)
	logoutUserCmd := appcommand.NewLogoutUserCommand(tokenBlacklistRepo)
	refreshTokenCmd := appcommand.NewRefreshTokenCommand(userRepo, jwtService, tokenBlacklistRepo)
	revokeTokenCmd := appcommand.NewRevokeTokenCommand(jwtService, tokenBlacklistRepo)
	updateProfileCmd := appcommand.NewUpdateProfileCommand(userRepo, userFactory, eventDispatcher)
	changePasswordCmd := appcommand.NewChangePasswordCommand(userRepo, eventDispatcher)

	// Setup queries
	getUserProfileQuery := appquery.NewGetUserProfileHandler(userRepo)
	getCurrentUserQuery := appquery.NewGetCurrentUserQuery(userRepo)
	validateTokenQuery := appquery.NewValidateTokenQuery(jwtService, tokenBlacklistRepo)

	logger.Component("auth.bootstrap").
		Info().
		Msg("auth service dependencies initialized")

	return &Dependencies{
		DB:                  db,
		Publisher:           publisher,
		Cache:               redisCache,
		UserRepo:            userRepo,
		RoleRepo:            roleRepo,
		PermissionRepo:       permissionRepo,
		UserRoleRepo:        userRoleRepo,
		RolePermissionRepo:  rolePermissionRepo,
		TokenBlacklistRepo:  tokenBlacklistRepo,
		EventDispatcher:     eventDispatcher,
		JwtService:          jwtService,
		RegisterUserCmd:     registerUserCmd,
		LoginUserCmd:        loginUserCmd,
		LogoutUserCmd:       logoutUserCmd,
		RefreshTokenCmd:     refreshTokenCmd,
		RevokeTokenCmd:      revokeTokenCmd,
		UpdateProfileCmd:    updateProfileCmd,
		ChangePasswordCmd:   changePasswordCmd,
		GetUserProfileQuery: getUserProfileQuery,
		GetCurrentUserQuery: getCurrentUserQuery,
		ValidateTokenQuery:  validateTokenQuery,
	}, nil
}

func setupPublisher() (*eventbuspublisher.KafkaPublisher, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	publisher, err := eventbuspublisher.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("auth.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		return nil, err
	}

	logger.Component("auth.bootstrap").
		Info().
		Msg("kafka publisher initialized")

	return publisher, nil
}

func setupEventDispatcher(publisher *eventbuspublisher.KafkaPublisher) *event_dispatcher.Dispatcher {
	dispatcher := event_dispatcher.NewDispatcher()

	// Create event broker adapter (abstraction over infrastructure)
	eventBrokerAdapter := eventbuspublisher.NewEventBrokerAdapter(publisher)

	// Register UserCreated handler
	userCreatedHandler := event_handler.NewUserCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("UserCreated", userCreatedHandler)
	logger.Component("auth.bootstrap").
		Info().
		Str("event_type", "UserCreated").
		Str("handler", "UserCreatedHandler").
		Msg("registered event handler")

	// Register UserProfileUpdated handler
	userProfileUpdatedHandler := event_handler.NewUserProfileUpdatedHandler()
	dispatcher.RegisterHandler("UserProfileUpdated", userProfileUpdatedHandler)
	logger.Component("auth.bootstrap").
		Info().
		Str("event_type", "UserProfileUpdated").
		Str("handler", "UserProfileUpdatedHandler").
		Msg("registered event handler")

	// Register UserPasswordChanged handler
	userPasswordChangedHandler := event_handler.NewUserPasswordChangedHandler()
	dispatcher.RegisterHandler("UserPasswordChanged", userPasswordChangedHandler)
	logger.Component("auth.bootstrap").
		Info().
		Str("event_type", "UserPasswordChanged").
		Str("handler", "UserPasswordChangedHandler").
		Msg("registered event handler")

	logger.Component("auth.bootstrap").
		Info().
		Int("total_handlers", 3).
		Msg("event dispatcher configured")

	return dispatcher
}

func setupDatabase() (*gorm.DB, error) {
	dsn := config.GetEnv("AUTH_DATABASE_DSN", "postgres://auth_user:auth_password@localhost:5433/auth_service?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Component("auth.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect database")
		return nil, err
	}

	// Configure connection pool for high concurrency
	sqlDB, err := db.DB()
	if err != nil {
		logger.Component("auth.bootstrap").
			Error().
			Err(err).
			Msg("failed to get underlying sql.DB")
		return nil, err
	}

	// Set connection pool settings
	maxOpenConns := config.GetEnvInt("AUTH_DB_MAX_OPEN_CONNS", 50)
	maxIdleConns := config.GetEnvInt("AUTH_DB_MAX_IDLE_CONNS", 10)
	connMaxLifetime := config.GetEnvInt("AUTH_DB_CONN_MAX_LIFETIME_MINUTES", 5)

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

	logger.Component("auth.bootstrap").
		Info().
		Int("max_open_conns", maxOpenConns).
		Int("max_idle_conns", maxIdleConns).
		Int("conn_max_lifetime_minutes", connMaxLifetime).
		Msg("database connected with connection pool configured")

	return db, nil
}

func setupCache() (cache.Cache, error) {
	addr := config.GetEnv("REDIS_ADDR", "localhost:6379")
	password := config.GetEnv("REDIS_PASSWORD", "")
	db := config.GetEnvInt("REDIS_DB", 0)

	redisCache, err := cache.NewRedisCache(addr, password, db, "auth.cache")
	if err != nil {
		return nil, err
	}

	logger.Component("auth.bootstrap").
		Info().
		Str("addr", addr).
		Int("db", db).
		Msg("redis cache initialized")

	return redisCache, nil
}

