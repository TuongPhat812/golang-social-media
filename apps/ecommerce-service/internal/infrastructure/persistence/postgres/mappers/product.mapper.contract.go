package mappers

import "golang-social-media/apps/ecommerce-service/internal/domain/product"

// ProductMapper defines the contract for mapping between domain Product and persistence models
type ProductMapper interface {
	ToModel(p product.Product) ProductModel
	ToDomain(model ProductModel) product.Product
	ToDomainList(models []ProductModel) []product.Product
	ToModelList(products []product.Product) []ProductModel
}


