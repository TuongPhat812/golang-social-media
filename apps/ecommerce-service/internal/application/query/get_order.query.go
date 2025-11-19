package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"
)

var _ contracts.GetOrderQuery = (*getOrderQuery)(nil)

type getOrderQuery struct {
	repo orders.Repository
	log  *zerolog.Logger
}

func NewGetOrderQuery(repo orders.Repository) contracts.GetOrderQuery {
	return &getOrderQuery{
		repo: repo,
		log:  logger.Component("ecommerce.query.get_order"),
	}
}

func (q *getOrderQuery) Execute(ctx context.Context, orderID string) (order.Order, error) {
	orderModel, err := q.repo.FindByID(ctx, orderID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("order_id", orderID).
			Msg("failed to get order")
		return order.Order{}, err
	}

	q.log.Info().
		Str("order_id", orderID).
		Msg("order retrieved")

	return orderModel, nil
}

