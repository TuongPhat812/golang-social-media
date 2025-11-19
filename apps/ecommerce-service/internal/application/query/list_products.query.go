package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/pkg/logger"
)

var _ contracts.ListProductsQuery = (*listProductsQuery)(nil)

type listProductsQuery struct {
	repo products.Repository
	log  *zerolog.Logger
}

func NewListProductsQuery(repo products.Repository) contracts.ListProductsQuery {
	return &listProductsQuery{
		repo: repo,
		log:  logger.Component("ecommerce.query.list_products"),
	}
}

func (q *listProductsQuery) Execute(ctx context.Context, req contracts.ListProductsQueryRequest) ([]product.Product, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	var status *product.Status
	if req.Status != "" {
		s := product.Status(req.Status)
		status = &s
	}

	products, err := q.repo.List(ctx, status, limit, req.Offset)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("status", req.Status).
			Int("limit", limit).
			Msg("failed to list products")
		return nil, err
	}

	q.log.Info().
		Int("count", len(products)).
		Str("status", req.Status).
		Msg("products retrieved")

	return products, nil
}

