# Mapper Implementation Summary

## Overview

Đã implement dedicated mapper package cho tất cả services theo DDD best practices. Mappers được tổ chức theo layers và responsibilities rõ ràng.

## Structure

### 1. Persistence Mappers (Domain ↔ Model)

**Location:** `internal/infrastructure/persistence/mappers/`

**Purpose:** Chuyển đổi giữa domain entities và persistence models

#### Auth Service
- `user.mapper.go` - Maps User domain ↔ memory representation

#### Chat Service
- `message.mapper.go` - Maps Message domain ↔ MessageModel

#### E-commerce Service
- `product.mapper.go` - Maps Product domain ↔ ProductModel
- `order.mapper.go` - Maps Order domain ↔ OrderModel + OrderItemModels

### 2. DTO Mappers (Domain ↔ DTO)

**Location:** `internal/interfaces/{rest|grpc}/mappers/`

**Purpose:** Chuyển đổi giữa domain entities và DTOs cho API layer

#### Auth Service
- `interfaces/rest/mappers/user.mapper.go` - Maps User domain ↔ HTTP DTOs

#### Chat Service
- `interfaces/grpc/mappers/message.mapper.go` - Maps Message domain ↔ gRPC DTOs

#### E-commerce Service
- `interfaces/grpc/mappers/product.mapper.go` - Maps Product domain ↔ gRPC DTOs
- `interfaces/grpc/mappers/order.mapper.go` - Maps Order domain ↔ gRPC DTOs

## Mapper Patterns

### 1. Persistence Mapper Pattern

```go
type ProductMapper struct{}

func NewProductMapper() *ProductMapper {
    return &ProductMapper{}
}

func (m *ProductMapper) ToModel(p product.Product) ProductModel { ... }
func (m *ProductMapper) ToDomain(model ProductModel) product.Product { ... }
func (m *ProductMapper) ToDomainList(models []ProductModel) []product.Product { ... }
```

### 2. DTO Mapper Pattern

```go
type ProductDTOMapper struct{}

func NewProductDTOMapper() *ProductDTOMapper {
    return &ProductDTOMapper{}
}

func (m *ProductDTOMapper) ToProduct(p product.Product) *ecommercev1.Product { ... }
func (m *ProductDTOMapper) ToProductList(products []product.Product) []*ecommercev1.Product { ... }
func (m *ProductDTOMapper) ToCreateProductResponse(p product.Product) *ecommercev1.CreateProductResponse { ... }
```

## Usage in Repositories

### Before (Mapper methods in model):
```go
func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
    model := ProductModelFromDomain(*p)  // Direct call
    // ...
    *p = model.ToDomain()
    return nil
}
```

### After (Dedicated mapper):
```go
type ProductRepository struct {
    db     *gorm.DB
    mapper *mappers.ProductMapper
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
    model := r.mapper.ToModel(*p)  // Use injected mapper
    // ...
    *p = r.mapper.ToDomain(model)
    return nil
}
```

## Usage in Handlers

### Before (Direct mapping):
```go
func (h *ProductHandler) CreateProduct(...) {
    // ...
    return &ecommercev1.CreateProductResponse{
        Product: &ecommercev1.Product{
            Id:    product.ID,
            Name:  product.Name,
            // ... manual mapping
        },
    }, nil
}
```

### After (DTO mapper):
```go
type ProductHandler struct {
    deps      *bootstrap.Dependencies
    dtoMapper *mappers.ProductDTOMapper
}

func (h *ProductHandler) CreateProduct(...) {
    // ...
    return h.dtoMapper.ToCreateProductResponse(product), nil
}
```

## Benefits

1. **Separation of Concerns**: Mapping logic tách biệt khỏi business logic
2. **Reusability**: Mappers có thể reuse ở nhiều nơi
3. **Testability**: Dễ test mapper riêng biệt
4. **Maintainability**: Dễ maintain và update mapping logic
5. **Consistency**: Consistent mapping pattern across all services
6. **Scalability**: Dễ mở rộng cho complex mappings

## Files Created

### Auth Service
- `apps/auth-service/internal/infrastructure/persistence/mappers/user.mapper.go`
- `apps/auth-service/internal/interfaces/rest/mappers/user.mapper.go`

### Chat Service
- `apps/chat-service/internal/infrastructure/persistence/mappers/message.mapper.go`
- `apps/chat-service/internal/interfaces/grpc/mappers/message.mapper.go`

### E-commerce Service
- `apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers/product.mapper.go`
- `apps/ecommerce-service/internal/infrastructure/persistence/postgres/mappers/order.mapper.go`
- `apps/ecommerce-service/internal/interfaces/grpc/mappers/product.mapper.go`
- `apps/ecommerce-service/internal/interfaces/grpc/mappers/order.mapper.go`

## Files Refactored

### Repositories
- `apps/chat-service/internal/infrastructure/persistence/message.repository.go`
- `apps/ecommerce-service/internal/infrastructure/persistence/postgres/product.repository.go`
- `apps/ecommerce-service/internal/infrastructure/persistence/postgres/order.repository.go`

### Handlers
- `apps/chat-service/internal/interfaces/grpc/chat/chat.handler.go`
- `apps/ecommerce-service/internal/interfaces/grpc/ecommerce/product.handler.go`
- `apps/ecommerce-service/internal/interfaces/grpc/ecommerce/order.handler.go`

## Next Steps

1. ✅ Remove old mapper methods from model files (optional cleanup)
2. ⏳ Create event mappers for Domain Events ↔ External Events
3. ⏳ Add mapper tests
4. ⏳ Document mapper patterns in code comments

## Notes

- Mappers are stateless and can be singleton instances
- Mappers should be pure functions (no side effects)
- Mappers handle nil and empty cases gracefully
- Batch mapping methods (`ToDomainList`, `ToProductList`) for performance

