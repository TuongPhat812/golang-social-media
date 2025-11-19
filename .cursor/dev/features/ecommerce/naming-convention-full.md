# E-commerce Service - Full Hậu tố Naming Convention

## Pattern hậu tố đầy đủ

Toàn bộ files trong ecommerce-service đã được đặt tên theo pattern hậu tố để phân biệt rõ ràng loại file.

## Domain Layer

```
domain/
├── order/
│   ├── order.entity.go              # Order Aggregate Root
│   ├── order_item.value_object.go  # OrderItem Value Object
│   └── order.event.go               # Order Domain Events
│
├── product/
│   ├── product.entity.go            # Product Entity
│   └── product.event.go             # Product Domain Events
│
├── shared/
│   ├── money.value_object.go        # Money Value Object
│   └── quantity.value_object.go    # Quantity Value Object
│
└── services/
    ├── stock_reservation.service.go
    ├── order_calculation.service.go
    └── product_availability.service.go
```

## Application Layer

```
application/
├── command/
│   ├── create_product.command.go
│   ├── create_order.command.go
│   ├── add_order_item.command.go
│   ├── confirm_order.command.go
│   ├── cancel_order.command.go
│   ├── update_product_stock.command.go
│   └── contracts/
│       ├── command.contract.go
│       ├── create_product.command.contract.go
│       ├── create_order.command.contract.go
│       ├── add_order_item.command.contract.go
│       ├── confirm_order.command.contract.go
│       ├── cancel_order.command.contract.go
│       └── update_product_stock.command.contract.go
│
├── query/
│   ├── get_product.query.go
│   ├── list_products.query.go
│   ├── get_order.query.go
│   ├── list_user_orders.query.go
│   └── contracts/
│       ├── query.contract.go
│       ├── get_product.query.contract.go
│       ├── list_products.query.contract.go
│       ├── get_order.query.contract.go
│       └── list_user_orders.query.contract.go
│
├── event_dispatcher/
│   └── dispatcher.service.go        # Event Dispatcher Service
│
├── event_handler/
│   ├── product_created.handler.go
│   ├── product_stock_updated.handler.go
│   ├── order_created.handler.go
│   ├── order_confirmed.handler.go
│   ├── order_cancelled.handler.go
│   └── contracts/
│       └── event_broker.contract.go
│
├── products/
│   └── repository.contract.go       # Product Repository Interface
│
└── orders/
    └── repository.contract.go       # Order Repository Interface
```

## Infrastructure Layer

```
infrastructure/
├── bootstrap/
│   └── bootstrap.service.go         # Bootstrap Service
│
├── grpc/
│   └── server.grpc.go               # gRPC Server
│
├── eventbus/
│   ├── publisher/
│   │   ├── kafka.publisher.go      # Kafka Publisher
│   │   ├── event_broker.adapter.go # Event Broker Adapter
│   │   └── contracts/
│   │       └── publisher.contract.go
│   │
│   └── subscriber/
│       ├── product_created.subscriber.go
│       ├── product_stock_updated.subscriber.go
│       ├── order_created.subscriber.go
│       ├── order_item_added.subscriber.go
│       ├── order_confirmed.subscriber.go
│       ├── order_cancelled.subscriber.go
│       └── contracts/
│           └── subscriber.contract.go
│
└── persistence/
    └── postgres/
        ├── product.model.go          # Product Model
        ├── product.repository.go   # Product Repository Implementation
        ├── order.model.go            # Order Model
        └── order.repository.go      # Order Repository Implementation
```

## Interfaces Layer

```
interfaces/
└── grpc/
    ├── register.service.go          # gRPC Service Registration
    └── ecommerce/
        ├── product_handler.go       # Product gRPC Handler
        └── order_handler.go         # Order gRPC Handler
```

## Quy tắc đặt tên (Full Pattern)

### Domain Layer:
- **`.entity.go`** → Entities (Order, Product)
- **`.value_object.go`** → Value Objects (OrderItem, Money, Quantity)
- **`.event.go`** → Domain Events (OrderCreatedEvent, ProductCreatedEvent)
- **`.service.go`** → Domain Services

### Application Layer:
- **`.command.go`** → Commands
- **`.command.contract.go`** → Command Contracts/Interfaces
- **`.query.go`** → Queries
- **`.query.contract.go`** → Query Contracts/Interfaces
- **`.handler.go`** → Event Handlers
- **`.contract.go`** → Contracts/Interfaces
- **`.service.go`** → Services (Event Dispatcher)
- **`.repository.contract.go`** → Repository Interfaces

### Infrastructure Layer:
- **`.model.go`** → Persistence Models
- **`.repository.go`** → Repository Implementations
- **`.publisher.go`** → Event Publishers
- **`.subscriber.go`** → Event Subscribers/Consumers
- **`.adapter.go`** → Adapters
- **`.contract.go`** → Contracts/Interfaces
- **`.service.go`** → Services (Bootstrap)
- **`.grpc.go`** → gRPC Server

### Interfaces Layer:
- **`.handler.go`** → gRPC Handlers
- **`.service.go`** → Service Registration

## Lợi ích

1. **Rõ ràng loại file**: Biết ngay đây là entity, command, query, handler, repository, publisher, subscriber, etc.
2. **Dễ tìm trong IDE**: Search `*.command.go`, `*.query.go`, `*.handler.go`, `*.repository.go`, `*.publisher.go`, `*.subscriber.go`
3. **Không conflict**: Khi có entity tên "Event" → `event.entity.go` vs `event.event.go`
4. **Consistent**: Pattern nhất quán trong toàn bộ codebase
5. **Dễ maintain**: Thay đổi một phần không ảnh hưởng phần khác
6. **Dễ test**: Test từng component riêng biệt
7. **Tuân thủ DDD**: Cấu trúc phản ánh đúng các khái niệm DDD

## Ví dụ: Tìm tất cả subscribers

```bash
find . -name "*.subscriber.go"
```

Kết quả:
```
infrastructure/eventbus/subscriber/product_created.subscriber.go
infrastructure/eventbus/subscriber/order_created.subscriber.go
...
```

## Ví dụ: Tìm tất cả publishers

```bash
find . -name "*.publisher.go"
```

Kết quả:
```
infrastructure/eventbus/publisher/kafka.publisher.go
```

## Kết luận

Full hậu tố pattern giúp:
- ✅ Phân biệt rõ ràng các loại file
- ✅ Dễ tìm và maintain
- ✅ Tránh conflict
- ✅ Consistent trong toàn bộ codebase
- ✅ Phù hợp với pattern Node.js/TypeScript
- ✅ Publisher và Subscriber đều có pattern rõ ràng
