# So sánh các cách tổ chức Mapper

## Cách 1: Mapper methods trong Model (hiện tại)

```go
// infrastructure/persistence/postgres/product.model.go
type ProductModel struct {
    ID    string `gorm:"column:id;type:uuid;primaryKey"`
    Name  string `gorm:"column:name;type:text;not null"`
    Price float64 `gorm:"column:price;type:decimal(10,2);not null"`
}

func ProductModelFromDomain(p product.Product) ProductModel {
    return ProductModel{
        ID:    p.ID,
        Name:  p.Name,
        Price: p.Price,
    }
}

func (m ProductModel) ToDomain() product.Product {
    return product.Product{
        ID:    m.ID,
        Name:  m.Name,
        Price: m.Price,
    }
}
```

### ✅ Ưu điểm:
1. **Đơn giản**: Mapping logic gần với model, dễ tìm
2. **Ít boilerplate**: Không cần tạo mapper struct
3. **Co-location**: Model và mapper cùng file, dễ maintain
4. **Phù hợp cho simple cases**: Khi mapping 1-1, không phức tạp

### ❌ Nhược điểm:
1. **Model phải import domain**: Tạo dependency từ infrastructure → domain
2. **Khó reuse**: Không thể dùng mapper ở nơi khác (ví dụ: DTO mapping)
3. **Khó test riêng**: Phải test cùng với model
4. **Khó scale**: Khi mapping phức tạp (nhiều Value Objects, nested structures)

---

## Cách 2: Dedicated Mapper Package

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
        ID:    p.ID,
        Name:  p.Name,
        Price: p.Price,
    }
}

func (m *ProductMapper) ToDomain(model postgres.ProductModel) product.Product {
    return product.Product{
        ID:    model.ID,
        Name:  model.Name,
        Price: model.Price,
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

### ✅ Ưu điểm:
1. **Separation of concerns**: Mapper logic tách biệt khỏi model
2. **Reusable**: Có thể dùng mapper cho nhiều mục đích (DTO, events)
3. **Testable**: Dễ test mapper riêng biệt
4. **Scalable**: Dễ mở rộng cho complex mappings
5. **Dependency direction đúng**: Infrastructure → Domain (không ngược lại)

### ❌ Nhược điểm:
1. **Nhiều files hơn**: Phải tạo thêm mapper package
2. **Boilerplate**: Phải tạo mapper struct cho mỗi entity
3. **Overkill cho simple cases**: Quá phức tạp khi mapping đơn giản

---

## So sánh chi tiết

| Tiêu chí | Mapper trong Model | Dedicated Mapper Package |
|----------|-------------------|-------------------------|
| **Đơn giản** | ✅ Rất đơn giản | ⚠️ Phức tạp hơn |
| **Co-location** | ✅ Cùng file với model | ❌ File riêng |
| **Reusability** | ❌ Chỉ dùng cho Model↔Domain | ✅ Có thể dùng cho nhiều mục đích |
| **Testability** | ⚠️ Test cùng model | ✅ Test riêng biệt |
| **Scalability** | ❌ Khó scale | ✅ Dễ scale |
| **Dependency** | ⚠️ Model import domain | ✅ Đúng direction |
| **Complex mappings** | ❌ Khó handle | ✅ Dễ handle |
| **Boilerplate** | ✅ Ít code | ❌ Nhiều code hơn |

---

## Khuyến nghị

### ✅ Dùng Mapper trong Model khi:

1. **Mapping đơn giản (1-1)**
   - Domain entity và persistence model gần giống nhau
   - Chỉ có primitive types, không có Value Objects
   - Không có nested structures phức tạp

2. **Project nhỏ/medium**
   - Ít entities
   - Không có nhiều mapping scenarios
   - Team nhỏ

3. **Ví dụ phù hợp:**
   ```go
   // Simple case - OK với mapper trong model
   type MessageModel struct {
       ID      string
       Content string
   }
   
   func (m MessageModel) ToDomain() domain.Message {
       return domain.Message{
           ID:      m.ID,
           Content: m.Content,
       }
   }
   ```

### ✅ Dùng Dedicated Mapper Package khi:

1. **Mapping phức tạp**
   - Có Value Objects (Money, Quantity, Address)
   - Có nested structures (Order với OrderItems)
   - Có nhiều transformations

2. **Cần reuse mapper**
   - Dùng cho nhiều mục đích: Model↔Domain, Domain↔DTO, Domain↔Event
   - Có nhiều output formats (JSON, gRPC, GraphQL)

3. **Project lớn/complex**
   - Nhiều entities
   - Nhiều mapping scenarios
   - Team lớn, cần maintainability

4. **Ví dụ cần dedicated mapper:**
   ```go
   // Complex case - Nên dùng dedicated mapper
   type Product struct {
       ID    string
       Price Money      // Value Object
       Stock Quantity   // Value Object
       Tags  []Tag      // Nested
   }
   
   type ProductModel struct {
       ID          string
       PriceAmount int64    // Flattened
       PriceCurrency string
       StockValue  int
       TagsJSON    string   // Serialized
   }
   ```

---

## Recommendation cho codebase hiện tại

### Phân tích codebase:

1. **Chat Service** - Simple mapping:
   ```go
   // Message: Primitive types only
   type Message struct {
       ID         string
       SenderID   string
       ReceiverID string
       Content    string
       CreatedAt  time.Time
   }
   ```
   ✅ **OK với mapper trong model** - Đơn giản, 1-1 mapping

2. **E-commerce Service** - Có Value Objects:
   ```go
   // Product: Có Value Objects (Money, Quantity)
   type Product struct {
       ID    string
       Price Money      // Value Object
       Stock Quantity   // Value Object
   }
   ```
   ⚠️ **Nên dùng dedicated mapper** - Có Value Objects, cần transformation

3. **Auth Service** - Simple:
   ```go
   // User: Primitive types
   type User struct {
       ID       string
       Email    string
       Password string
       Name     string
   }
   ```
   ✅ **OK với mapper trong model** - Đơn giản

### Kết luận:

**Hybrid approach** - Dùng cả 2 cách tùy complexity:

1. **Simple entities** (Message, User): Giữ mapper trong model
2. **Complex entities** (Product, Order với Value Objects): Chuyển sang dedicated mapper

---

## Migration path

### Bước 1: Giữ nguyên simple cases
```go
// ✅ Giữ nguyên
func (m MessageModel) ToDomain() domain.Message { ... }
```

### Bước 2: Tạo dedicated mapper cho complex cases
```go
// ✅ Tạo mới cho Product
// infrastructure/persistence/mappers/product.mapper.go
type ProductMapper struct{}

func (m *ProductMapper) ToDomain(model ProductModel) product.Product {
    return product.Product{
        ID:    model.ID,
        Price: money.New(model.PriceAmount, model.PriceCurrency),
        Stock: quantity.New(model.StockValue),
    }
}
```

### Bước 3: Dần migrate khi cần
- Khi entity trở nên phức tạp → migrate sang dedicated mapper
- Khi cần reuse mapper → migrate sang dedicated mapper

---

## Best Practice Summary

1. **Start simple**: Bắt đầu với mapper trong model
2. **Migrate when needed**: Chuyển sang dedicated mapper khi:
   - Mapping trở nên phức tạp
   - Cần reuse mapper
   - Có Value Objects hoặc nested structures
3. **Be consistent**: Trong cùng một service, dùng cùng một pattern
4. **Document**: Ghi rõ khi nào dùng cách nào

