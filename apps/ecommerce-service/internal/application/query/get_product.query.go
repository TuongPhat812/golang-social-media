package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/products"
	"golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	"golang-social-media/pkg/logger"
)

var _ contracts.GetProductQuery = (*getProductQuery)(nil)

type getProductQuery struct {
	repo products.Repository
	log  *zerolog.Logger
}

func NewGetProductQuery(repo products.Repository) contracts.GetProductQuery {
	return &getProductQuery{
		repo: repo,
		log:  logger.Component("ecommerce.query.get_product"),
	}
}

func (q *getProductQuery) Execute(ctx context.Context, productID string) (product.Product, error) {
	productModel, err := q.repo.FindByID(ctx, productID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("product_id", productID).
			Msg("failed to get product")
		return product.Product{}, err
	}

	q.log.Info().
		Str("product_id", productID).
		Msg("product retrieved")

	return productModel, nil
}

