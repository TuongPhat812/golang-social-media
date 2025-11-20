package mappers

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductDTOMapperImpl implements ProductDTOMapper interface
type ProductDTOMapperImpl struct{}

var _ ProductDTOMapper = (*ProductDTOMapperImpl)(nil)

// NewProductDTOMapper creates a new ProductDTOMapperImpl
func NewProductDTOMapper() ProductDTOMapper {
	return &ProductDTOMapperImpl{}
}

// ToCreateProductRequest converts gRPC CreateProductRequest to command request
// Note: This extracts data for command, not domain entity directly
func (m *ProductDTOMapperImpl) FromCreateProductRequest(req *ecommercev1.CreateProductRequest) (name, description string, price float64, stock int) {
	return req.GetName(), req.GetDescription(), req.GetPrice(), int(req.GetStock())
}

// ToProduct converts domain Product to gRPC Product
func (m *ProductDTOMapperImpl) ToProduct(p product.Product) *ecommercev1.Product {
	return &ecommercev1.Product{
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

// ToProductList converts a slice of domain Products to gRPC Products
func (m *ProductDTOMapperImpl) ToProductList(products []product.Product) []*ecommercev1.Product {
	result := make([]*ecommercev1.Product, len(products))
	for i, p := range products {
		result[i] = m.ToProduct(p)
	}
	return result
}

// ToCreateProductResponse converts domain Product to CreateProductResponse
func (m *ProductDTOMapperImpl) ToCreateProductResponse(p product.Product) *ecommercev1.CreateProductResponse {
	return &ecommercev1.CreateProductResponse{
		Product: m.ToProduct(p),
	}
}

// ToGetProductResponse converts domain Product to GetProductResponse
func (m *ProductDTOMapperImpl) ToGetProductResponse(p product.Product) *ecommercev1.GetProductResponse {
	return &ecommercev1.GetProductResponse{
		Product: m.ToProduct(p),
	}
}

// ToListProductsResponse converts domain Products to ListProductsResponse
func (m *ProductDTOMapperImpl) ToListProductsResponse(products []product.Product) *ecommercev1.ListProductsResponse {
	return &ecommercev1.ListProductsResponse{
		Products: m.ToProductList(products),
	}
}

// ToUpdateProductStockResponse converts domain Product to UpdateProductStockResponse
func (m *ProductDTOMapperImpl) ToUpdateProductStockResponse(p product.Product) *ecommercev1.UpdateProductStockResponse {
	return &ecommercev1.UpdateProductStockResponse{
		Product: m.ToProduct(p),
	}
}

