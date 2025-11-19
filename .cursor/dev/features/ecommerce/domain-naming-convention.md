# Domain Naming Convention - Phân biệt Entity và Domain Events

## Vấn đề

Khi có domain entity tên là **"Event"** (ví dụ: Event trong hệ thống quản lý sự kiện), làm sao phân biệt:
- **Event Entity** (business model) - Event có ID, Name, Date, Location
- **Domain Events** (EventCreatedEvent, EventUpdatedEvent) - các sự kiện business

## Giải pháp: Naming Convention

### Cấu trúc folder:

```
domain/event/                    # Event Entity (business model)
├── event.go                     # Event Entity (business model)
├── event_item.go                # EventItem Value Object (nếu có)
└── domain_events.go             # Domain Events (EventCreatedEvent, etc.)
```

### Hoặc pattern khác:

```
domain/event/                    # Event Entity (business model)
├── event.go                     # Event Entity (business model)
├── event_item.go                # EventItem Value Object
└── events.go                    # Domain Events (có thể dùng nếu không conflict)
```

## Quy tắc đặt tên

### 1. Entity file: Tên entity (lowercase)
- `event.go` → Event entity (business model)
- `order.go` → Order entity
- `product.go` → Product entity

### 2. Domain Events file: `domain_events.go` hoặc `events.go`
- `domain_events.go` → Rõ ràng hơn, tránh nhầm lẫn
- `events.go` → Ngắn gọn, OK nếu không có conflict

### 3. Value Objects: Tên value object
- `order_item.go` → OrderItem value object
- `event_item.go` → EventItem value object

## Ví dụ cụ thể: Event Management System

### Cấu trúc:

```
domain/event/
├── event.go                     # Event Entity (business model)
│   └── type Event struct {
│       ID          string
│       Name        string
│       Date        time.Time
│       Location    string
│       ...
│   }
│
├── event_item.go                # EventItem Value Object (nếu có)
│   └── type EventItem struct {
│       ...
│   }
│
└── domain_events.go             # Domain Events (EventCreatedEvent, etc.)
    ├── type EventCreatedEvent struct { ... }
    ├── type EventUpdatedEvent struct { ... }
    └── type EventCancelledEvent struct { ... }
```

### Code example:

```go
// domain/event/event.go
package event

type Event struct {
    ID          string
    Name        string
    Date        time.Time
    Location    string
    // ...
    events []DomainEvent  // Internal domain events
}

func (e *Event) Create() {
    e.addEvent(EventCreatedEvent{
        EventID: e.ID,
        Name:    e.Name,
        // ...
    })
}

// domain/event/domain_events.go
package event

// DomainEvent interface
type DomainEvent interface {
    Type() string
}

// EventCreatedEvent - Domain event khi Event entity được tạo
type EventCreatedEvent struct {
    EventID string
    Name    string
    Date    string
    // ...
}

func (e EventCreatedEvent) Type() string {
    return "EventCreated"  // Domain event type
}
```

## So sánh các pattern

### Pattern 1: `domain_events.go` (Khuyến nghị)
```
domain/event/
├── event.go           → Event Entity
└── domain_events.go   → Domain Events (rõ ràng, không nhầm lẫn)
```

**Ưu điểm:**
- ✅ Rõ ràng, không nhầm lẫn
- ✅ Dễ phân biệt Entity vs Domain Events
- ✅ Phù hợp khi có entity tên "Event"

### Pattern 2: `events.go` (OK nếu không conflict)
```
domain/event/
├── event.go    → Event Entity
└── events.go   → Domain Events
```

**Ưu điểm:**
- ✅ Ngắn gọn
- ⚠️ Có thể nhầm lẫn nếu không đọc kỹ

### Pattern 3: Tên cụ thể hơn
```
domain/event/
├── event.go                    → Event Entity
└── event_domain_events.go      → Domain Events (quá dài)
```

**Nhược điểm:**
- ❌ Tên quá dài
- ❌ Không cần thiết

## Best Practice

### 1. Nếu có entity tên "Event":
```
domain/event/
├── event.go           → Event Entity
└── domain_events.go   → Domain Events ✅ (rõ ràng nhất)
```

### 2. Nếu không có entity tên "Event":
```
domain/order/
├── order.go    → Order Entity
└── events.go  → Domain Events ✅ (OK, không conflict)
```

### 3. Nếu có nhiều loại events:
```
domain/event/
├── event.go                    → Event Entity
├── domain_events.go            → Domain Events (EventCreated, EventUpdated)
└── event_notification_events.go → Notification Events (nếu cần tách riêng)
```

## Kết luận

**Quy tắc chung:**
1. **Entity file**: Tên entity (lowercase) - `event.go`, `order.go`
2. **Domain Events file**: 
   - `domain_events.go` → Khi có entity tên "Event" (rõ ràng nhất)
   - `events.go` → Khi không có conflict (ngắn gọn)
3. **Value Objects**: Tên value object - `order_item.go`, `event_item.go`

**Khuyến nghị:**
- Dùng `domain_events.go` khi có entity tên "Event"
- Dùng `events.go` khi không có conflict
- Luôn comment rõ ràng trong code

