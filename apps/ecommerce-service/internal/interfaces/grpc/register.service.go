package grpc

import (
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	ecommercegrpc "golang-social-media/apps/ecommerce-service/internal/interfaces/grpc/ecommerce"
	grpcmappers "golang-social-media/apps/ecommerce-service/internal/interfaces/grpc/mappers"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/grpc"
)

// RegisterServices registers all gRPC services
func RegisterServices(server *grpc.Server, deps *bootstrap.Dependencies) {
	// Setup DTO mappers
	productDTOMapper := grpcmappers.NewProductDTOMapper()
	orderDTOMapper := grpcmappers.NewOrderDTOMapper()

	// Register Product Service
	productHandler := ecommercegrpc.NewProductHandler(deps, productDTOMapper)
	ecommercev1.RegisterProductServiceServer(server, productHandler)

	// Register Order Service
	orderHandler := ecommercegrpc.NewOrderHandler(deps, orderDTOMapper)
	ecommercev1.RegisterOrderServiceServer(server, orderHandler)
}

