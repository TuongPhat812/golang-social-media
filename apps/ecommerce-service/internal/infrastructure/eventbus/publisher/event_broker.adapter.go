package publisher

import (
	"context"

	appcontracts "golang-social-media/apps/ecommerce-service/internal/application/event_handler/contracts"
	infracontracts "golang-social-media/apps/ecommerce-service/internal/infrastructure/eventbus/publisher/contracts"
)

var _ appcontracts.EventBrokerPublisher = (*EventBrokerAdapter)(nil)

// EventBrokerAdapter adapts KafkaPublisher to EventBrokerPublisher interface
type EventBrokerAdapter struct {
	publisher infracontracts.EcommercePublisher
}

func NewEventBrokerAdapter(publisher infracontracts.EcommercePublisher) appcontracts.EventBrokerPublisher {
	return &EventBrokerAdapter{
		publisher: publisher,
	}
}

func (a *EventBrokerAdapter) PublishProductCreated(ctx context.Context, payload appcontracts.ProductCreatedPayload) error {
	return a.publisher.PublishProductCreated(ctx, infracontracts.ProductCreated{
		ProductID:   payload.ProductID,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		CreatedAt:   payload.CreatedAt,
	})
}

func (a *EventBrokerAdapter) PublishProductStockUpdated(ctx context.Context, payload appcontracts.ProductStockUpdatedPayload) error {
	return a.publisher.PublishProductStockUpdated(ctx, infracontracts.ProductStockUpdated{
		ProductID: payload.ProductID,
		OldStock:  payload.OldStock,
		NewStock:  payload.NewStock,
		UpdatedAt: payload.UpdatedAt,
	})
}

func (a *EventBrokerAdapter) PublishOrderCreated(ctx context.Context, payload appcontracts.OrderCreatedPayload) error {
	return a.publisher.PublishOrderCreated(ctx, infracontracts.OrderCreated{
		OrderID:     payload.OrderID,
		UserID:      payload.UserID,
		TotalAmount: payload.TotalAmount,
		ItemCount:   payload.ItemCount,
		CreatedAt:   payload.CreatedAt,
	})
}

func (a *EventBrokerAdapter) PublishOrderItemAdded(ctx context.Context, payload appcontracts.OrderItemAddedPayload) error {
	return a.publisher.PublishOrderItemAdded(ctx, infracontracts.OrderItemAdded{
		OrderID:   payload.OrderID,
		ProductID: payload.ProductID,
		Quantity:  payload.Quantity,
		UnitPrice: payload.UnitPrice,
		SubTotal:  payload.SubTotal,
		UpdatedAt: payload.UpdatedAt,
	})
}

func (a *EventBrokerAdapter) PublishOrderConfirmed(ctx context.Context, payload appcontracts.OrderConfirmedPayload) error {
	return a.publisher.PublishOrderConfirmed(ctx, infracontracts.OrderConfirmed{
		OrderID:     payload.OrderID,
		UserID:      payload.UserID,
		TotalAmount: payload.TotalAmount,
		ItemCount:   payload.ItemCount,
		ConfirmedAt: payload.ConfirmedAt,
	})
}

func (a *EventBrokerAdapter) PublishOrderCancelled(ctx context.Context, payload appcontracts.OrderCancelledPayload) error {
	return a.publisher.PublishOrderCancelled(ctx, infracontracts.OrderCancelled{
		OrderID:    payload.OrderID,
		UserID:     payload.UserID,
		CancelledAt: payload.CancelledAt,
	})
}

