# Mappers trong Domain-Driven Design

## Khái niệm

**Mapper** (hay **Domain Mapper**) là một pattern trong DDD dùng để chuyển đổi giữa các representations khác nhau của cùng một data trong các layer khác nhau. Mapper giúp:

1. **Tách biệt concerns**: Domain layer không biết về persistence details
2. **Bảo vệ domain model**: Tránh leak infrastructure concerns vào domain
3. **Linh hoạt**: Dễ dàng thay đổi persistence layer mà không ảnh hưởng domain
4. **Testability**: Dễ test domain logic độc lập với persistence

## Các loại Mapper trong DDD

### 1. Domain Entity ↔ Persistence Model

Chuyển đổi giữa:
- **Domain Entity**: Pure business logic, không có annotations của ORM
- **Persistence Model**: Có GORM tags, database constraints, indexes

**Ví dụ:**
```go
// Domain Entity (pure business logic)
type Product struct {
    ID          string
    Name        string
    Price       Money
    Stock       Quantity
    Status      Status
}

// Persistence Model (có GORM tags)
type ProductModel struct {
    ID          string    `gorm:"column:id;type:uuid;primaryKey"`
    Name        string    `gorm:"column:name;type:varchar(255);not null"`
    PriceAmount int64     `gorm:"column:price_amount;type:bigint;not null"`
    PriceCurrency string  `gorm:"column:price_currency;type:varchar(3);not null"`
    StockValue   int      `gorm:"column:stock_value;type:int;not null"`
    Status       string   `gorm:"column:status;type:varchar(20);not null"`
    CreatedAt    time.Time `gorm:"column:created_at;not null"`
    UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}
```

### 2. Domain Entity ↔ DTO/Contract

Chuyển đổi giữa:
- **Domain Entity**: Internal representation với business logic
- **DTO (Data Transfer Object)**: External representation cho API

**Ví dụ:**
```go
// Domain Entity
type User struct {
    ID       string
    Email    string
    Password string // Internal, không expose ra ngoài
    Name     string
}

// DTO (cho HTTP response)
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
    // Password không có trong DTO
}
```

### 3. Domain Event ↔ External Event

Chuyển đổi giữa:
- **Domain Event**: Internal event với domain-specific structure
- **External Event**: Event cho event bus (Kafka, etc.)

**Ví dụ:**
```go
// Domain Event
type UserCreatedEvent struct {
    UserID    string
    Email     string
    Name      string
    CreatedAt string
}

// External Event (cho Kafka)
type UserCreatedKafkaEvent struct {
    EventType string    `json:"event_type"`
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}
```

## Cách tổ chức Mapper trong codebase

### Option 1: Mapper methods trong Model (hiện tại)

```go
// infrastructure/persistence/postgres/product.model.go
func ProductModelFromDomain(p product.Product) ProductModel {
    return ProductModel{
        ID:          p.ID,
        Name:        p.Name,
        PriceAmount: p.Price.Amount(),
        // ...
    }
}

func (m ProductModel) ToDomain() product.Product {
    return product.Product{
        ID:    m.ID,
        Name:  m.Name,
        Price: money.New(m.PriceAmount, m.PriceCurrency),
        // ...
    }
}
```

**Ưu điểm:**
- Đơn giản, dễ hiểu
- Mapping logic gần với model

**Nhược điểm:**
- Model phải biết về domain structure
- Khó reuse khi có nhiều mapping scenarios

### Option 2: Dedicated Mapper package (recommended)

```
internal/
├── domain/
│   └── product/
│       └── product.entity.go
├── infrastructure/
│   └── persistence/
│       ├── postgres/
│       │   ├── product.model.go
│       │   └── product.repository.go
│       └── mappers/
│           ├── product.mapper.go      # Domain ↔ Model
│           └── product_dto.mapper.go # Domain ↔ DTO
└── interfaces/
    └── rest/
        └── mappers/
            └── product.mapper.go     # Domain ↔ HTTP DTO
```

**Ví dụ:**

```go
// infrastructure/persistence/mappers/product.mapper.go
package mappers

import (
    "golang-social-media/apps/ecommerce-service/internal/domain/product"
    "golang-social-media/apps/ecommerce-service/internal/infrastructure/persistence/postgres"
)

type ProductMapper struct{}

func NewProductMapper() *ProductMapper {
    return &ProductMapper{}
}

func (m *ProductMapper) ToModel(p product.Product) postgres.ProductModel {
    return postgres.ProductModel{
        ID:          p.ID,
        Name:        p.Name,
        PriceAmount: p.Price.Amount(),
        PriceCurrency: p.Price.Currency(),
        StockValue:  p.Stock.Value(),
        Status:      string(p.Status),
        CreatedAt:   p.CreatedAt,
        UpdatedAt:   p.UpdatedAt,
    }
}

func (m *ProductMapper) ToDomain(model postgres.ProductModel) product.Product {
    return product.Product{
        ID:     model.ID,
        Name:   model.Name,
        Price:  money.New(model.PriceAmount, model.PriceCurrency),
        Stock:  quantity.New(model.StockValue),
        Status: product.Status(model.Status),
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }
}

func (m *ProductMapper) ToDomainList(models []postgres.ProductModel) []product.Product {
    products := make([]product.Product, len(models))
    for i, model := range models {
        products[i] = m.ToDomain(model)
    }
    return products
}
```

```go
// interfaces/rest/mappers/product.mapper.go
package mappers

import (
    "golang-social-media/apps/ecommerce-service/internal/domain/product"
    "golang-social-media/pkg/contracts/ecommerce"
)

type ProductDTOMapper struct{}

func NewProductDTOMapper() *ProductDTOMapper {
    return &ProductDTOMapper{}
}

func (m *ProductDTOMapper) ToDTO(p product.Product) ecommerce.ProductResponse {
    return ecommerce.ProductResponse{
        ID:          p.ID,
        Name:        p.Name,
        Price:       p.Price.Amount(),
        Currency:    p.Price.Currency(),
        Stock:       p.Stock.Value(),
        Status:      string(p.Status),
        CreatedAt:   p.CreatedAt,
    }
}

func (m *ProductDTOMapper) ToDTOList(products []product.Product) []ecommerce.ProductResponse {
    dtos := make([]ecommerce.ProductResponse, len(products))
    for i, p := range products {
        dtos[i] = m.ToDTO(p)
    }
    return dtos
}
```

## Khi nào cần Mapper?

### ✅ Cần Mapper khi:

1. **Domain Entity khác Persistence Model**
   - Domain dùng Value Objects (Money, Quantity)
   - Persistence dùng primitive types (int64, string)
   - Domain không có timestamps
   - Persistence cần timestamps cho audit

2. **Domain Entity khác DTO**
   - Domain có sensitive fields (password, tokens)
   - DTO cần format khác (dates, numbers)
   - DTO có thêm metadata (pagination, links)

3. **Domain Event khác External Event**
   - Domain event có internal structure
   - External event cần schema cho event bus
   - Cần thêm metadata (version, source)

### ❌ Không cần Mapper khi:

1. **Domain Entity = Persistence Model** (simple cases)
   - Chỉ có primitive types
   - Không có Value Objects
   - Không có sensitive data

2. **Domain Entity = DTO** (internal APIs)
   - Không có security concerns
   - Không cần format khác

## Best Practices

### 1. Mapper nên là pure functions

```go
// ✅ Good: Pure function
func ToDomain(model ProductModel) product.Product {
    return product.Product{...}
}

// ❌ Bad: Side effects
func (m *ProductMapper) ToDomain(model ProductModel) product.Product {
    m.stats.Increment() // Side effect
    return product.Product{...}
}
```

### 2. Handle nil và empty cases

```go
func (m *ProductMapper) ToDomain(model *ProductModel) (*product.Product, error) {
    if model == nil {
        return nil, errors.NewNotFoundError(errors.CodeProductNotFound)
    }
    return &product.Product{...}, nil
}
```

### 3. Batch mapping cho performance

```go
func (m *ProductMapper) ToDomainList(models []ProductModel) []product.Product {
    products := make([]product.Product, 0, len(models))
    for _, model := range models {
        products = append(products, m.ToDomain(model))
    }
    return products
}
```

### 4. Validate trong mapper

```go
func (m *ProductMapper) ToDomain(model ProductModel) (product.Product, error) {
    if model.ID == "" {
        return product.Product{}, errors.New("invalid model: missing ID")
    }
    // ...
}
```

### 5. Mapper cho complex aggregates

```go
// Order có nhiều OrderItems
func (m *OrderMapper) ToDomain(orderModel OrderModel, itemModels []OrderItemModel) order.Order {
    items := make([]order.OrderItem, len(itemModels))
    for i, itemModel := range itemModels {
        items[i] = m.itemMapper.ToDomain(itemModel)
    }
    
    return order.Order{
        ID:    orderModel.ID,
        Items: items,
        // ...
    }
}
```

## Ví dụ trong codebase hiện tại

### Chat Service

```go
// infrastructure/persistence/message.model.go
func MessageModelFromDomain(msg domain.Message) MessageModel {
    return MessageModel{
        ID:         msg.ID,
        SenderID:   msg.SenderID,
        ReceiverID: msg.ReceiverID,
        Content:    msg.Content,
        CreatedAt:  msg.CreatedAt,
    }
}

func (m MessageModel) ToDomain() domain.Message {
    return domain.Message{
        ID:         m.ID,
        SenderID:   m.SenderID,
        ReceiverID: m.ReceiverID,
        Content:    m.Content,
        CreatedAt:  m.CreatedAt,
    }
}
```

### E-commerce Service

```go
// infrastructure/persistence/postgres/product.model.go
func ProductModelFromDomain(p product.Product) ProductModel {
    return ProductModel{
        ID:          p.ID,
        Name:        p.Name,
        PriceAmount: p.Price.Amount(),
        PriceCurrency: p.Price.Currency(),
        StockValue:  p.Stock.Value(),
        Status:      string(p.Status),
    }
}

func (m ProductModel) ToDomain() product.Product {
    return product.Product{
        ID:    m.ID,
        Name:  m.Name,
        Price: money.New(m.PriceAmount, m.PriceCurrency),
        Stock: quantity.New(m.StockValue),
        Status: product.Status(m.Status),
    }
}
```

## Tóm tắt

Mapper trong DDD là cầu nối giữa các layers:
- **Domain ↔ Infrastructure**: Bảo vệ domain khỏi persistence details
- **Domain ↔ Interfaces**: Bảo vệ domain khỏi external concerns
- **Domain Events ↔ External Events**: Chuyển đổi events cho event bus

Mapper giúp:
- ✅ Tách biệt concerns
- ✅ Bảo vệ domain model
- ✅ Dễ test và maintain
- ✅ Linh hoạt khi thay đổi infrastructure

