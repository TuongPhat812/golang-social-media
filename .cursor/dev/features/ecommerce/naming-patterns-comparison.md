# Naming Patterns Comparison

## Pattern 1: Tên file dài (hiện tại)

```
domain/order/
├── order.go                    # Order Entity
├── order_item.go               # OrderItem Value Object
└── events.go                   # Domain Events
```

**Ưu điểm:**
- ✅ Tên file ngắn gọn
- ✅ Dễ đọc trong file explorer
- ✅ Phù hợp với Go convention

**Nhược điểm:**
- ❌ Khó phân biệt loại file khi có nhiều file
- ❌ Khi có entity tên "Event" → conflict với `events.go`

## Pattern 2: Hậu tố như Node.js (user đề xuất)

```
domain/order/
├── order.entity.go             # Order Entity
├── order_item.value_object.go # OrderItem Value Object
└── order.event.go              # Domain Events
```

**Ưu điểm:**
- ✅ Rõ ràng loại file ngay từ tên
- ✅ Dễ phân biệt: `.entity.go`, `.event.go`, `.service.go`
- ✅ Không conflict khi có entity tên "Event"
- ✅ Consistent với pattern Node.js
- ✅ Dễ tìm trong IDE (search "*.entity.go")

**Nhược điểm:**
- ⚠️ Tên file dài hơn một chút
- ⚠️ Không phải Go convention truyền thống

## Pattern 3: Hybrid (cân bằng)

```
domain/order/
├── order.go                    # Order Entity (ngắn gọn)
├── order_item.go               # OrderItem Value Object
└── order.events.go             # Domain Events (có hậu tố)
```

## So sánh chi tiết

### Trường hợp Order Domain:

**Pattern 1 (hiện tại):**
```
order/
├── order.go
├── order_item.go
└── events.go
```

**Pattern 2 (hậu tố):**
```
order/
├── order.entity.go
├── order_item.value_object.go
└── order.event.go
```

**Pattern 3 (hybrid):**
```
order/
├── order.go
├── order_item.go
└── order.events.go
```

### Trường hợp Event Domain (có conflict):

**Pattern 1 (hiện tại):**
```
event/
├── event.go
└── domain_events.go  ← Phải đổi tên để tránh conflict
```

**Pattern 2 (hậu tố):**
```
event/
├── event.entity.go
└── event.event.go    ← Rõ ràng, không conflict
```

**Pattern 3 (hybrid):**
```
event/
├── event.go
└── event.events.go   ← Rõ ràng, không conflict
```

## Khuyến nghị

### Option 1: Full hậu tố (rõ ràng nhất)
```
domain/order/
├── order.entity.go
├── order_item.value_object.go
└── order.event.go

domain/event/
├── event.entity.go
└── event.event.go
```

### Option 2: Hybrid (cân bằng)
```
domain/order/
├── order.go
├── order_item.go
└── order.events.go

domain/event/
├── event.go
└── event.events.go
```

### Option 3: Minimal hậu tố (chỉ khi cần)
```
domain/order/
├── order.go
├── order_item.go
└── events.go  ← OK nếu không conflict

domain/event/
├── event.go
└── event.events.go  ← Dùng hậu tố khi có conflict
```

## Quy tắc đặt tên với hậu tố

1. **Entity**: `{name}.entity.go` hoặc `{name}.go`
2. **Value Object**: `{name}.value_object.go` hoặc `{name}.go`
3. **Domain Events**: `{name}.event.go` hoặc `{name}.events.go`
4. **Domain Service**: `{name}.service.go`

## Kết luận

**Pattern hậu tố có ưu điểm:**
- ✅ Rõ ràng loại file
- ✅ Không conflict
- ✅ Dễ tìm trong IDE
- ✅ Consistent

**Nếu user thích pattern hậu tố, tôi có thể refactor lại!**

