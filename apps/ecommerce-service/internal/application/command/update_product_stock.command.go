package command

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/ecommerce-service/internal/application/event_dispatcher"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/pkg/logger"
)

var _ contracts.UpdateProductStockCommand = (*updateProductStockCommand)(nil)

type updateProductStockCommand struct {
	repo           products.Repository
	eventDispatcher *event_dispatcher.Dispatcher
	log            *zerolog.Logger
}

func NewUpdateProductStockCommand(
	repo products.Repository,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.UpdateProductStockCommand {
	return &updateProductStockCommand{
		repo:           repo,
		eventDispatcher: eventDispatcher,
		log:            logger.Component("ecommerce.command.update_product_stock"),
	}
}

func (c *updateProductStockCommand) Execute(ctx context.Context, productID string, newStock int) error {
	// Load product
	productModel, err := c.repo.FindByID(ctx, productID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("product_id", productID).
			Msg("failed to find product")
		return err
	}

	// Update stock (domain logic - adds domain event)
	if err := productModel.UpdateStock(newStock); err != nil {
		c.log.Error().
			Err(err).
			Str("product_id", productID).
			Int("new_stock", newStock).
			Msg("failed to update stock")
		return err
	}

	// Save domain events BEFORE persisting
	domainEvents := productModel.Events()

	// Persist updated product
	if err := c.repo.Update(ctx, &productModel); err != nil {
		c.log.Error().
			Err(err).
			Str("product_id", productID).
			Msg("failed to update product")
		return err
	}

	// Dispatch domain events AFTER successful persistence
	productModel.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("product_id", productID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("product_id", productID).
		Int("new_stock", newStock).
		Msg("product stock updated")

	return nil
}

