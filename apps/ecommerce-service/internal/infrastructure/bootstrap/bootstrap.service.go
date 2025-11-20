package bootstrap

import (
	"context"

	appcommand "golang-social-media/apps/ecommerce-service/internal/application/command"
	commandcontracts "golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	event_handler "golang-social-media/apps/ecommerce-service/internal/application/event_handler"
	apporders "golang-social-media/apps/ecommerce-service/internal/application/orders"
	appproducts "golang-social-media/apps/ecommerce-service/internal/application/products"
	appquery "golang-social-media/apps/ecommerce-service/internal/application/query"
	querycontracts "golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	unit_of_work "golang-social-media/apps/ecommerce-service/internal/application/unit_of_work"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/cache"
	eventbuspublisher "golang-social-media/apps/ecommerce-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/eventstore"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/outbox"
	postgrespersistence "golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
	postgrespersistencemappers "golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers"
	domainfactories "golang-social-media/apps/ecommerce-service/internal/domain/factories"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	DB                    *gorm.DB
	Publisher             *eventbuspublisher.KafkaPublisher
	ProductRepo           appproducts.Repository
	OrderRepo             apporders.Repository
	UoWFactory            unit_of_work.Factory
	EventDispatcher       *event_dispatcher.Dispatcher
	OutboxService         *outbox.OutboxService
	OutboxProcessor       *outbox.Processor
	EventStore            *eventstore.EventStoreService
	Cache                 cache.Cache
	BatchRepo             *postgrespersistence.BatchRepository
	CreateProductCmd      commandcontracts.CreateProductCommand
	UpdateProductStockCmd commandcontracts.UpdateProductStockCommand
	CreateOrderCmd        commandcontracts.CreateOrderCommand
	AddOrderItemCmd       commandcontracts.AddOrderItemCommand
	ConfirmOrderCmd       commandcontracts.ConfirmOrderCommand
	CancelOrderCmd        commandcontracts.CancelOrderCommand
	GetProductQuery       querycontracts.GetProductQuery
	ListProductsQuery     querycontracts.ListProductsQuery
	GetOrderQuery         querycontracts.GetOrderQuery
	ListUserOrdersQuery   querycontracts.ListUserOrdersQuery
}

// SetupDependencies initializes all service dependencies
func SetupDependencies(ctx context.Context) (*Dependencies, error) {
	// Setup database
	db, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	// Setup Redis cache
	redisCache, err := setupCache()
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Warn().
			Err(err).
			Msg("failed to setup cache, continuing without cache")
		redisCache = nil
	}

	// Setup cache wrappers
	var productCache *cache.ProductCache
	var orderCache *cache.OrderCache
	if redisCache != nil {
		productCache = cache.NewProductCache(redisCache)
		orderCache = cache.NewOrderCache(redisCache)
	}

	// Setup mappers
	productMapper := postgrespersistencemappers.NewProductMapper()
	orderMapper := postgrespersistencemappers.NewOrderMapper()

	// Setup repositories with cache and mappers
	productRepo := postgrespersistence.NewProductRepository(db, productMapper, productCache)
	orderRepo := postgrespersistence.NewOrderRepository(db, orderMapper, orderCache)

	// Setup batch repository
	batchRepo := postgrespersistence.NewBatchRepository(db)

	// Setup Unit of Work factory with mappers
	uowFactory := postgrespersistence.NewUnitOfWorkFactory(db, productMapper, orderMapper)

	// Setup event bus publisher
	publisher, err := setupPublisher()
	if err != nil {
		return nil, err
	}

	// Setup event dispatcher (for backward compatibility)
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup Outbox
	outboxRepo := postgrespersistence.NewOutboxRepository(db)
	outboxService := outbox.NewOutboxService(outboxRepo)
	outboxProcessor := outbox.NewProcessor(outboxRepo, publisher)

	// Setup Event Store
	eventStoreRepo := postgrespersistence.NewEventStoreRepository(db)
	eventStoreService := eventstore.NewEventStoreService(eventStoreRepo)

	// Setup factories
	orderFactory := domainfactories.NewOrderFactory()

	// Setup commands
	commands := setupCommands(uowFactory, orderFactory, productRepo, orderRepo, eventDispatcher, outboxService, eventStoreService)

	// Setup queries
	queries := setupQueries(productRepo, orderRepo)

	logger.Component("ecommerce.bootstrap").
		Info().
		Msg("ecommerce service dependencies initialized")

	return &Dependencies{
		DB:                    db,
		Publisher:             publisher,
		ProductRepo:           productRepo,
		OrderRepo:             orderRepo,
		UoWFactory:            uowFactory,
		EventDispatcher:       eventDispatcher,
		OutboxService:         outboxService,
		OutboxProcessor:       outboxProcessor,
		EventStore:            eventStoreService,
		Cache:                 redisCache,
		BatchRepo:             batchRepo,
		CreateProductCmd:       commands.CreateProduct,
		UpdateProductStockCmd: commands.UpdateProductStock,
		CreateOrderCmd:        commands.CreateOrder,
		AddOrderItemCmd:       commands.AddOrderItem,
		ConfirmOrderCmd:       commands.ConfirmOrder,
		CancelOrderCmd:        commands.CancelOrder,
		GetProductQuery:       queries.GetProduct,
		ListProductsQuery:     queries.ListProducts,
		GetOrderQuery:         queries.GetOrder,
		ListUserOrdersQuery:   queries.ListUserOrders,
	}, nil
}

func setupDatabase() (*gorm.DB, error) {
	dsn := config.GetEnv("ECOMMERCE_DATABASE_DSN", "postgres://ecommerce_user:ecommerce_password@localhost:5432/ecommerce_service?sslmode=disable")
	
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect database")
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to get underlying sql.DB")
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)                 // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)                 // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(5 * 60 * 1000000000) // 5 minutes

	logger.Component("ecommerce.bootstrap").
		Info().
		Int("max_open_conns", 25).
		Int("max_idle_conns", 10).
		Msg("database connected with connection pooling")

	return db, nil
}

func setupPublisher() (*eventbuspublisher.KafkaPublisher, error) {
	brokers := config.GetEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	publisher, err := eventbuspublisher.NewKafkaPublisher(brokers)
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to create kafka publisher")
		return nil, err
	}

	logger.Component("ecommerce.bootstrap").
		Info().
		Msg("kafka publisher initialized")

	return publisher, nil
}

func setupCache() (cache.Cache, error) {
	addr := config.GetEnv("REDIS_ADDR", "localhost:6379")
	password := config.GetEnv("REDIS_PASSWORD", "")
	db := config.GetEnvInt("REDIS_DB", 0)

	redisCache, err := cache.NewRedisCache(addr, password, db)
	if err != nil {
		return nil, err
	}

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("addr", addr).
		Int("db", db).
		Msg("redis cache initialized")

	return redisCache, nil
}

func setupEventDispatcher(publisher *eventbuspublisher.KafkaPublisher) *event_dispatcher.Dispatcher {
	dispatcher := event_dispatcher.NewDispatcher()

	// Create event broker adapter (abstraction over infrastructure)
	eventBrokerAdapter := eventbuspublisher.NewEventBrokerAdapter(publisher)

	// Register ProductCreated handler
	productCreatedHandler := event_handler.NewProductCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("ProductCreated", productCreatedHandler)
	logger.Component("ecommerce.bootstrap").
		Info().
		Str("event_type", "ProductCreated").
		Str("handler", "ProductCreatedHandler").
		Msg("registered event handler")

	// Register ProductStockUpdated handler
	productStockUpdatedHandler := event_handler.NewProductStockUpdatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("ProductStockUpdated", productStockUpdatedHandler)
	logger.Component("ecommerce.bootstrap").
		Info().
		Str("event_type", "ProductStockUpdated").
		Str("handler", "ProductStockUpdatedHandler").
		Msg("registered event handler")

	// Register OrderCreated handler
	orderCreatedHandler := event_handler.NewOrderCreatedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("OrderCreated", orderCreatedHandler)
	logger.Component("ecommerce.bootstrap").
		Info().
		Str("event_type", "OrderCreated").
		Str("handler", "OrderCreatedHandler").
		Msg("registered event handler")

	// Register OrderConfirmed handler
	orderConfirmedHandler := event_handler.NewOrderConfirmedHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("OrderConfirmed", orderConfirmedHandler)
	logger.Component("ecommerce.bootstrap").
		Info().
		Str("event_type", "OrderConfirmed").
		Str("handler", "OrderConfirmedHandler").
		Msg("registered event handler")

	// Register OrderCancelled handler
	orderCancelledHandler := event_handler.NewOrderCancelledHandler(eventBrokerAdapter)
	dispatcher.RegisterHandler("OrderCancelled", orderCancelledHandler)
	logger.Component("ecommerce.bootstrap").
		Info().
		Str("event_type", "OrderCancelled").
		Str("handler", "OrderCancelledHandler").
		Msg("registered event handler")

	logger.Component("ecommerce.bootstrap").
		Info().
		Int("total_handlers", 5).
		Msg("event dispatcher configured")

	return dispatcher
}

type commands struct {
	CreateProduct      commandcontracts.CreateProductCommand
	UpdateProductStock commandcontracts.UpdateProductStockCommand
	CreateOrder        commandcontracts.CreateOrderCommand
	AddOrderItem       commandcontracts.AddOrderItemCommand
	ConfirmOrder       commandcontracts.ConfirmOrderCommand
	CancelOrder        commandcontracts.CancelOrderCommand
}

func setupCommands(
	uowFactory unit_of_work.Factory,
	orderFactory domainfactories.OrderFactory,
	productRepo appproducts.Repository,
	orderRepo apporders.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
	outboxService *outbox.OutboxService,
	eventStore *eventstore.EventStoreService,
) commands {
	createProductCmd := appcommand.NewCreateProductCommand(productRepo, eventDispatcher)
	updateProductStockCmd := appcommand.NewUpdateProductStockCommand(productRepo, eventDispatcher)
	createOrderCmd := appcommand.NewCreateOrderCommand(uowFactory, orderFactory, outboxService, eventStore)
	addOrderItemCmd := appcommand.NewAddOrderItemCommand(uowFactory, eventDispatcher)
	confirmOrderCmd := appcommand.NewConfirmOrderCommand(uowFactory, eventDispatcher)
	cancelOrderCmd := appcommand.NewCancelOrderCommand(uowFactory, eventDispatcher)

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "CreateProductCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "UpdateProductStockCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "CreateOrderCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "AddOrderItemCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "ConfirmOrderCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("command", "CancelOrderCommand").
		Msg("registered command")

	logger.Component("ecommerce.bootstrap").
		Info().
		Int("total_commands", 6).
		Msg("commands configured")

	return commands{
		CreateProduct:      createProductCmd,
		UpdateProductStock: updateProductStockCmd,
		CreateOrder:        createOrderCmd,
		AddOrderItem:       addOrderItemCmd,
		ConfirmOrder:        confirmOrderCmd,
		CancelOrder:         cancelOrderCmd,
	}
}

type queries struct {
	GetProduct     querycontracts.GetProductQuery
	ListProducts   querycontracts.ListProductsQuery
	GetOrder       querycontracts.GetOrderQuery
	ListUserOrders querycontracts.ListUserOrdersQuery
}

func setupQueries(
	productRepo appproducts.Repository,
	orderRepo apporders.Repository,
) queries {
	getProductQuery := appquery.NewGetProductQuery(productRepo)
	listProductsQuery := appquery.NewListProductsQuery(productRepo)
	getOrderQuery := appquery.NewGetOrderQuery(orderRepo)
	listUserOrdersQuery := appquery.NewListUserOrdersQuery(orderRepo)

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("query", "GetProductQuery").
		Msg("registered query")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("query", "ListProductsQuery").
		Msg("registered query")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("query", "GetOrderQuery").
		Msg("registered query")

	logger.Component("ecommerce.bootstrap").
		Info().
		Str("query", "ListUserOrdersQuery").
		Msg("registered query")

	logger.Component("ecommerce.bootstrap").
		Info().
		Int("total_queries", 4).
		Msg("queries configured")

	return queries{
		GetProduct:     getProductQuery,
		ListProducts:   listProductsQuery,
		GetOrder:       getOrderQuery,
		ListUserOrders: listUserOrdersQuery,
	}
}

