package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/pkg/logger"
)

var _ contracts.CreateProductCommand = (*createProductCommand)(nil)

type createProductCommand struct {
	repo           products.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log            *zerolog.Logger
}

func NewCreateProductCommand(
	repo products.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.CreateProductCommand {
	return &createProductCommand{
		repo:           repo,
		eventDispatcher: eventDispatcher,
		log:            logger.Component("ecommerce.command.create_product"),
	}
}

func (c *createProductCommand) Execute(ctx context.Context, req contracts.CreateProductCommandRequest) (product.Product, error) {
	createdAt := time.Now().UTC()
	productModel := product.Product{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      product.StatusActive,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	// Set status based on stock
	if productModel.Stock == 0 {
		productModel.Status = product.StatusOutOfStock
	}

	// Validate business rules before persisting or publishing
	if err := productModel.Validate(); err != nil {
		c.log.Error().
			Err(err).
			Str("name", req.Name).
			Msg("product validation failed")
		return product.Product{}, err
	}

	// Domain logic: create product (this adds domain events internally)
	productModel.Create()

	// Save domain events BEFORE persisting (repository might overwrite the product)
	domainEvents := productModel.Events()

	// Persist to database
	if err := c.repo.Create(ctx, &productModel); err != nil {
		c.log.Error().
			Err(err).
			Str("name", req.Name).
			Msg("failed to persist product")
		return product.Product{}, err
	}

	// Dispatch domain events AFTER successful persistence
	productModel.ClearEvents() // Clear events after dispatch

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("product_id", productModel.ID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("product_id", productModel.ID).
		Str("name", productModel.Name).
		Msg("product created")

	return productModel, nil
}

