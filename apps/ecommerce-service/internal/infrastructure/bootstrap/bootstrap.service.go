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
	eventbuspublisher "golang-social-media/apps/ecommerce-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/config"
	"golang-social-media/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	DB                    *gorm.DB
	Publisher             *eventbuspublisher.KafkaPublisher
	ProductRepo           appproducts.Repository
	OrderRepo             apporders.Repository
	EventDispatcher       *event_dispatcher.Dispatcher
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

	// Setup repositories
	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	// Setup event bus publisher
	publisher, err := setupPublisher()
	if err != nil {
		return nil, err
	}

	// Setup event dispatcher
	eventDispatcher := setupEventDispatcher(publisher)

	// Setup commands
	commands := setupCommands(productRepo, orderRepo, eventDispatcher)

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
		EventDispatcher:       eventDispatcher,
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Component("ecommerce.bootstrap").
			Error().
			Err(err).
			Msg("failed to connect database")
		return nil, err
	}

	logger.Component("ecommerce.bootstrap").
		Info().
		Msg("database connected")

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
	productRepo appproducts.Repository,
	orderRepo apporders.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) commands {
	createProductCmd := appcommand.NewCreateProductCommand(productRepo, eventDispatcher)
	updateProductStockCmd := appcommand.NewUpdateProductStockCommand(productRepo, eventDispatcher)
	createOrderCmd := appcommand.NewCreateOrderCommand(orderRepo, eventDispatcher)
	addOrderItemCmd := appcommand.NewAddOrderItemCommand(orderRepo, productRepo, eventDispatcher)
	confirmOrderCmd := appcommand.NewConfirmOrderCommand(orderRepo, productRepo, eventDispatcher)
	cancelOrderCmd := appcommand.NewCancelOrderCommand(orderRepo, productRepo, eventDispatcher)

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

