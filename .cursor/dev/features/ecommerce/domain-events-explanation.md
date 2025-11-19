# Domain Events - Giải thích chi tiết

## Domain Events là gì?

**Domain Events** là các sự kiện business xảy ra trong domain layer. Chúng **KHÔNG phải** là "event của event", mà là các **event types** (loại sự kiện).

### Ví dụ:
- `OrderCreated` - Sự kiện: Order được tạo
- `OrderConfirmed` - Sự kiện: Order được xác nhận
- `ProductStockUpdated` - Sự kiện: Stock của product được cập nhật

## Cấu trúc file trong thực tế

### Cấu trúc chuẩn cho Order Aggregate:

```
domain/order/
├── order.go          # Order Aggregate Root (entity chính)
├── order_item.go     # OrderItem Value Object
└── events.go         # Domain Events (các event types)
```

### Giải thích từng file:

#### 1. `order.go` - Aggregate Root
```go
// Order represents an order aggregate root
type Order struct {
    ID          string
    UserID      string
    Status      Status
    Items       []OrderItem  // Value Objects
    // ...
    events []DomainEvent     // Internal events list
}

// Business methods emit domain events
func (o *Order) Create() {
    o.addEvent(OrderCreatedEvent{...})  // Emit event
}
```

#### 2. `order_item.go` - Value Object
```go
// OrderItem represents an item in an order (Value Object)
type OrderItem struct {
    ProductID string
    Quantity  int
    // ...
}
```

#### 3. `events.go` - Domain Events (KHÔNG phải "event của event")
```go
// DomainEvent interface
type DomainEvent interface {
    Type() string
}

// OrderCreatedEvent - Event type: Order được tạo
type OrderCreatedEvent struct {
    OrderID     string
    UserID      string
    // ...
}

func (e OrderCreatedEvent) Type() string {
    return "OrderCreated"  // Event type name
}
```

## Domain Events KHÔNG phải là "event của event"

### ❌ SAI: "Event của event"
```
OrderCreatedEvent
  └── Event của OrderCreatedEvent  ❌ KHÔNG CÓ
```

### ✅ ĐÚNG: Domain Events là các event types
```
OrderCreatedEvent      → Event type: Order được tạo
OrderConfirmedEvent    → Event type: Order được xác nhận
OrderCancelledEvent    → Event type: Order bị hủy
```

## Flow của Domain Events

### 1. Domain Entity emit events
```go
// order.go
func (o *Order) Create() {
    // Business logic
    o.Status = StatusDraft
    
    // Emit domain event
    o.addEvent(OrderCreatedEvent{
        OrderID: o.ID,
        UserID: o.UserID,
        // ...
    })
}
```

### 2. Application layer dispatch events
```go
// command/create_order.command.go
func (c *createOrderCommand) Execute(...) {
    // Create order (emits events internally)
    orderModel.Create()
    
    // Get events from domain
    domainEvents := orderModel.Events()
    
    // Dispatch to handlers
    for _, event := range domainEvents {
        c.eventDispatcher.Dispatch(ctx, event)
    }
}
```

### 3. Event Handlers transform to infrastructure events
```go
// event_handler/order_created.handler.go
func (h *OrderCreatedHandler) Handle(ctx, domainEvent) {
    // Transform domain event → infrastructure event
    payload := contracts.OrderCreatedPayload{
        OrderID: domainEvent.OrderID,
        // ...
    }
    
    // Publish to Kafka
    h.eventBroker.PublishOrderCreated(ctx, payload)
}
```

## Tại sao tách file?

### ✅ Cấu trúc rõ ràng:
```
order/
├── order.go          → Aggregate Root (entity chính)
├── order_item.go     → Value Object (phụ thuộc)
└── events.go         → Domain Events (event types)
```

### ❌ KHÔNG nên gộp tất cả vào 1 file:
```
order.go  → Order + OrderItem + Events  ❌ Quá dài, khó đọc
```

## So sánh với các service khác

### Chat Service:
```
domain/message/
├── entity.go    → Message entity
└── event.go     → MessageCreatedEvent
```

### Notification Service:
```
domain/notification/
├── entity.go    → Notification entity
└── event.go     → NotificationCreatedEvent, NotificationReadEvent
```

### E-commerce Service (chuẩn):
```
domain/order/
├── order.go         → Order Aggregate Root
├── order_item.go    → OrderItem Value Object
└── events.go        → OrderCreatedEvent, OrderConfirmedEvent, etc.
```

## Kết luận

1. **Domain Events** = Các event types (OrderCreated, OrderConfirmed, etc.)
2. **KHÔNG có** "event của event"
3. **Cấu trúc file**:
   - `order.go` - Aggregate Root
   - `order_item.go` - Value Object
   - `events.go` - Domain Events (các event types)
4. **Mỗi file một mục đích** - Dễ đọc, dễ maintain

