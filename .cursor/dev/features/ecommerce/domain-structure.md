# E-commerce Service - Domain Layer Structure

## Cấu trúc Domain Layer (Full Hậu tố Pattern)

Domain layer được tổ chức rõ ràng với pattern hậu tố để phân biệt các khái niệm DDD:

```
domain/
├── product/                              # Product Aggregate
│   ├── product.entity.go                # Product Entity
│   └── product.event.go                 # Product Domain Events
│
├── order/                                # Order Aggregate
│   ├── order.entity.go                  # Order Aggregate Root
│   ├── order_item.value_object.go       # OrderItem Value Object
│   └── order.event.go                    # Order Domain Events
│
├── shared/                               # Shared Value Objects
│   ├── money.value_object.go             # Money Value Object
│   └── quantity.value_object.go         # Quantity Value Object
│
└── services/                             # Domain Services
    ├── stock_reservation.service.go
    ├── order_calculation.service.go
    └── product_availability.service.go
```

## Phân loại theo hậu tố

### 1. **`.entity.go`** - Entities
- **`order/order.entity.go`** - `Order` Aggregate Root
  - Quản lý OrderItems (aggregate boundary)
  - Entry point để access Order aggregate
  - Có business logic: Create, AddItem, Confirm, Cancel

- **`product/product.entity.go`** - `Product` Entity
  - Standalone entity, không phải aggregate root
  - Có identity và business logic

### 2. **`.value_object.go`** - Value Objects
- **`order/order_item.value_object.go`** - `OrderItem` value object
  - Immutable, defined by values
  - Factory method: `NewOrderItem()`
  
- **`shared/money.value_object.go`** - `Money` value object
  - Amount + Currency
  - Immutable, type-safe
  
- **`shared/quantity.value_object.go`** - `Quantity` value object
  - Quantity với validation
  - Immutable, type-safe

### 3. **`.event.go`** - Domain Events
- **`product/product.event.go`** - Product domain events
  - ProductCreated
  - ProductStockUpdated

- **`order/order.event.go`** - Order domain events
  - OrderCreated
  - OrderItemAdded
  - OrderConfirmed
  - OrderCancelled

### 4. **`.service.go`** - Domain Services
- **`services/stock_reservation.service.go`** - Stock reservation logic
- **`services/order_calculation.service.go`** - Order calculation logic
- **`services/product_availability.service.go`** - Product availability logic

## Quy tắc đặt tên (Full Hậu tố Pattern)

1. **Entity**: `{name}.entity.go` (ví dụ: `order.entity.go`, `product.entity.go`)
2. **Value Object**: `{name}.value_object.go` (ví dụ: `order_item.value_object.go`, `money.value_object.go`)
3. **Domain Events**: `{name}.event.go` (ví dụ: `order.event.go`, `product.event.go`)
4. **Domain Service**: `{name}.service.go` (ví dụ: `order_calculation.service.go`)

## Lợi ích của Full Hậu tố Pattern

1. **Rõ ràng loại file**: Biết ngay đây là entity, value object, event hay service
2. **Dễ tìm trong IDE**: Search `*.entity.go`, `*.event.go`, `*.service.go`
3. **Không conflict**: Khi có entity tên "Event" → `event.entity.go` vs `event.event.go`
4. **Consistent**: Pattern nhất quán trong toàn bộ domain layer
5. **Dễ maintain**: Thay đổi một phần không ảnh hưởng phần khác
6. **Dễ test**: Test từng component riêng biệt
7. **Tuân thủ DDD**: Cấu trúc phản ánh đúng các khái niệm DDD

## Ví dụ: Trường hợp có Entity tên "Event"

```
domain/event/
├── event.entity.go        # Event Entity (business model)
└── event.event.go         # Domain Events (EventCreatedEvent, etc.)
```

**Rõ ràng, không nhầm lẫn!**

## So sánh với các service khác

### Chat Service (có thể refactor sau):
```
domain/message/
├── entity.go              # Message Entity
└── event.go               # MessageCreatedEvent
```

### E-commerce Service (Full hậu tố):
```
domain/order/
├── order.entity.go        # Order Aggregate Root
├── order_item.value_object.go # OrderItem Value Object
└── order.event.go         # Order Domain Events
```

## Kết luận

Full hậu tố pattern giúp:
- ✅ Phân biệt rõ ràng các loại file
- ✅ Dễ tìm và maintain
- ✅ Tránh conflict
- ✅ Consistent trong toàn bộ codebase
