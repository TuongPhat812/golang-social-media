package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/orders"
	"golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/order"
	"golang-social-media/pkg/logger"
)

var _ contracts.ListUserOrdersQuery = (*listUserOrdersQuery)(nil)

type listUserOrdersQuery struct {
	repo orders.Repository
	log  *zerolog.Logger
}

func NewListUserOrdersQuery(repo orders.Repository) contracts.ListUserOrdersQuery {
	return &listUserOrdersQuery{
		repo: repo,
		log:  logger.Component("ecommerce.query.list_user_orders"),
	}
}

func (q *listUserOrdersQuery) Execute(ctx context.Context, userID string, limit int) ([]order.Order, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	orders, err := q.repo.ListByUser(ctx, userID, limit)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Int("limit", limit).
			Msg("failed to list user orders")
		return nil, err
	}

	q.log.Info().
		Str("user_id", userID).
		Int("count", len(orders)).
		Msg("user orders retrieved")

	return orders, nil
}

