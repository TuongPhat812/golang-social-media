package grpc

import (
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/bootstrap"
	ecommercegrpc "golang-social-media/apps/ecommerce-service/internal/interfaces/grpc/ecommerce"
	ecommercev1 "golang-social-media/pkg/gen/ecommerce/v1"
	"google.golang.org/grpc"
)

// RegisterServices registers all gRPC services
func RegisterServices(server *grpc.Server, deps *bootstrap.Dependencies) {
	// Register Product Service
	productHandler := ecommercegrpc.NewProductHandler(deps)
	ecommercev1.RegisterProductServiceServer(server, productHandler)

	// Register Order Service
	orderHandler := ecommercegrpc.NewOrderHandler(deps)
	ecommercev1.RegisterOrderServiceServer(server, orderHandler)
}

