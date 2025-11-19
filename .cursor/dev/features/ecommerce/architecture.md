# E-commerce Module Architecture

## Domain Model

### Entities:
1. **Product** - Sản phẩm
   - ID, Name, Description, Price, Stock, Status
   - Business logic: Validate(), UpdateStock(), IsAvailable()

2. **Order** - Đơn hàng
   - ID, UserID, Status, TotalAmount, CreatedAt, UpdatedAt
   - Business logic: Create(), AddItem(), CalculateTotal(), Confirm(), Cancel()
   - Aggregate Root cho OrderItems

3. **OrderItem** - Item trong đơn hàng
   - OrderID, ProductID, Quantity, UnitPrice, SubTotal
   - Value Object hoặc Entity (phụ thuộc vào business rules)

### Domain Events:
- `ProductCreated`
- `ProductStockUpdated`
- `OrderCreated`
- `OrderConfirmed`
- `OrderCancelled`
- `OrderItemAdded`

## Cấu trúc File

```
apps/ecommerce-service/
├── cmd/
│   └── ecommerce-service/
│       └── main.go
├── go.mod
├── go.sum
└── internal/
    ├── application/
    │   ├── command/
    │   │   ├── contracts/
    │   │   │   ├── command.contract.go
    │   │   │   ├── create_product.command.contract.go
    │   │   │   ├── update_product_stock.command.contract.go
    │   │   │   ├── create_order.command.contract.go
    │   │   │   ├── add_order_item.command.contract.go
    │   │   │   ├── confirm_order.command.contract.go
    │   │   │   └── cancel_order.command.contract.go
    │   │   ├── dto/
    │   │   │   ├── create_product.command.dto.go
    │   │   │   ├── create_order.command.dto.go
    │   │   │   └── add_order_item.command.dto.go
    │   │   ├── create_product.command.go
    │   │   ├── update_product_stock.command.go
    │   │   ├── create_order.command.go
    │   │   ├── add_order_item.command.go
    │   │   ├── confirm_order.command.go
    │   │   └── cancel_order.command.go
    │   ├── query/
    │   │   ├── contracts/
    │   │   │   ├── query.contract.go
    │   │   │   ├── get_product.query.contract.go
    │   │   │   ├── list_products.query.contract.go
    │   │   │   ├── get_order.query.contract.go
    │   │   │   └── list_user_orders.query.contract.go
    │   │   ├── dto/
    │   │   │   ├── product.query.dto.go
    │   │   │   └── order.query.dto.go
    │   │   ├── get_product.query.go
    │   │   ├── list_products.query.go
    │   │   ├── get_order.query.go
    │   │   └── list_user_orders.query.go
    │   ├── event_dispatcher/
    │   │   └── dispatcher.go
    │   └── event_handler/
    │       ├── contracts/
    │       │   └── event_broker.contract.go
    │       ├── product_created.handler.go
    │       ├── product_stock_updated.handler.go
    │       ├── order_created.handler.go
    │       ├── order_confirmed.handler.go
    │       └── order_cancelled.handler.go
    ├── domain/
    │   ├── product/
    │   │   ├── entity.go
    │   │   └── event.go
    │   ├── order/
    │   │   ├── entity.go
    │   │   ├── order_item.go (Value Object hoặc Entity)
    │   │   └── event.go
    │   └── shared/
    │       └── value_objects.go (Money, Quantity, etc.)
    ├── infrastructure/
    │   ├── bootstrap/
    │   │   └── bootstrap.go
    │   ├── eventbus/
    │   │   └── publisher/
    │   │       ├── contracts/
    │   │       │   ├── publisher.contract.go
    │   │       │   └── ecommerce_publisher.contract.go
    │   │       ├── kafka_publisher.go
    │   │       └── event_broker_adapter.go
    │   ├── persistence/
    │   │   └── postgres/
    │   │       ├── product_repository.db.go
    │   │       ├── order_repository.db.go
    │   │       ├── product_model.go
    │   │       └── order_model.go
    │   └── grpc/
    │       └── server.go
    └── interfaces/
        └── grpc/
            ├── register.go
            └── ecommerce/
                └── handler.go
```

## Use Cases

### Commands:
1. **CreateProduct** - Tạo sản phẩm mới
2. **UpdateProductStock** - Cập nhật số lượng tồn kho
3. **CreateOrder** - Tạo đơn hàng mới
4. **AddOrderItem** - Thêm item vào đơn hàng
5. **ConfirmOrder** - Xác nhận đơn hàng (trừ stock)
6. **CancelOrder** - Hủy đơn hàng (hoàn lại stock)

### Queries:
1. **GetProduct** - Lấy thông tin sản phẩm
2. **ListProducts** - Danh sách sản phẩm (có filter, pagination)
3. **GetOrder** - Lấy thông tin đơn hàng
4. **ListUserOrders** - Danh sách đơn hàng của user

## Database Schema

### Products Table:
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL, -- 'active', 'inactive', 'out_of_stock'
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Orders Table:
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL, -- 'draft', 'confirmed', 'cancelled', 'completed'
    total_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Order Items Table:
```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    sub_total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
```

## Domain Events Flow

### Product Created:
```
CreateProductCommand
  → Product.Create()
  → ProductCreatedEvent
  → EventDispatcher
  → ProductCreatedHandler
  → Kafka Publisher
```

### Order Confirmed:
```
ConfirmOrderCommand
  → Order.Confirm()
  → Validate stock availability
  → Update product stock
  → OrderConfirmedEvent
  → ProductStockUpdatedEvent
  → EventDispatcher
  → Handlers publish to Kafka
```

