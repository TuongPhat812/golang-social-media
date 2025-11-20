package mappers

import (
	"golang-social-media/apps/ecommerce-service/internal/domain/product"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
)

// ProductDTOMapper defines the contract for mapping between domain Product and gRPC DTOs
type ProductDTOMapper interface {
	FromCreateProductRequest(req *ecommercev1.CreateProductRequest) (name, description string, price float64, stock int)
	ToProduct(p product.Product) *ecommercev1.Product
	ToProductList(products []product.Product) []*ecommercev1.Product
	ToCreateProductResponse(p product.Product) *ecommercev1.CreateProductResponse
	ToGetProductResponse(p product.Product) *ecommercev1.GetProductResponse
	ToListProductsResponse(products []product.Product) *ecommercev1.ListProductsResponse
	ToUpdateProductStockResponse(p product.Product) *ecommercev1.UpdateProductStockResponse
}


