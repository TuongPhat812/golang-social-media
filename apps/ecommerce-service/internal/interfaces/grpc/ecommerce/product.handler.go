package ecommerce

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/ecommerce-service/internal/application/command/contracts"
	querycontracts "golang-social-media/apps/ecommerce-service/internal/application/query/contracts"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	"golang-social-media/pkg/logger"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductHandler struct {
	ecommercev1.UnimplementedProductServiceServer
	deps *bootstrap.Dependencies
	log  *zerolog.Logger
}

func NewProductHandler(deps *bootstrap.Dependencies) *ProductHandler {
	return &ProductHandler{
		deps: deps,
		log:  logger.Component("ecommerce.grpc.product"),
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *ecommercev1.CreateProductRequest) (*ecommercev1.CreateProductResponse, error) {
	cmdReq := contracts.CreateProductCommandRequest{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
	}

	product, err := h.deps.CreateProductCmd.Execute(ctx, cmdReq)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("name", req.GetName()).
			Msg("failed to create product")
		return nil, err
	}

	h.log.Info().
		Str("product_id", product.ID).
		Str("name", product.Name).
		Msg("product created")

	return &ecommercev1.CreateProductResponse{
		Product: &ecommercev1.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
			Status:      string(product.Status),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *ecommercev1.GetProductRequest) (*ecommercev1.GetProductResponse, error) {
	product, err := h.deps.GetProductQuery.Execute(ctx, req.GetProductId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("product_id", req.GetProductId()).
			Msg("failed to get product")
		return nil, err
	}

	return &ecommercev1.GetProductResponse{
		Product: &ecommercev1.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
			Status:      string(product.Status),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *ecommercev1.ListProductsRequest) (*ecommercev1.ListProductsResponse, error) {
	queryReq := querycontracts.ListProductsQueryRequest{
		Status: req.GetStatus(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}

	products, err := h.deps.ListProductsQuery.Execute(ctx, queryReq)
	if err != nil {
		h.log.Error().
			Err(err).
			Msg("failed to list products")
		return nil, err
	}

	pbProducts := make([]*ecommercev1.Product, len(products))
	for i, p := range products {
		pbProducts[i] = &ecommercev1.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
			Status:      string(p.Status),
			CreatedAt:   timestamppb.New(p.CreatedAt),
			UpdatedAt:   timestamppb.New(p.UpdatedAt),
		}
	}

	return &ecommercev1.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (h *ProductHandler) UpdateProductStock(ctx context.Context, req *ecommercev1.UpdateProductStockRequest) (*ecommercev1.UpdateProductStockResponse, error) {
	err := h.deps.UpdateProductStockCmd.Execute(ctx, req.GetProductId(), int(req.GetNewStock()))
	if err != nil {
		h.log.Error().
			Err(err).
			Str("product_id", req.GetProductId()).
			Int32("new_stock", req.GetNewStock()).
			Msg("failed to update product stock")
		return nil, err
	}

	// Fetch updated product
	product, err := h.deps.GetProductQuery.Execute(ctx, req.GetProductId())
	if err != nil {
		h.log.Error().
			Err(err).
			Str("product_id", req.GetProductId()).
			Msg("failed to get updated product")
		return nil, err
	}

	return &ecommercev1.UpdateProductStockResponse{
		Product: &ecommercev1.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
			Status:      string(product.Status),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
	}, nil
}

