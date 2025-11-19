# Ví dụ: Event Domain với Event Entity

## Tình huống

Có domain entity tên **"Event"** (ví dụ: Event trong hệ thống quản lý sự kiện - Event Management System).

Làm sao phân biệt:
- **Event Entity** (business model) - Event có ID, Name, Date, Location
- **Domain Events** (EventCreatedEvent, EventUpdatedEvent) - các sự kiện business

## Cấu trúc folder

```
domain/event/
├── event.go           # Event Entity (business model)
└── domain_events.go   # Domain Events (EventCreatedEvent, etc.)
```

## Code Example

### 1. `event.go` - Event Entity (Business Model)

```go
package event

import (
    "errors"
    "time"
)

// Event represents an event entity (business model)
// This is a business entity, NOT a domain event
type Event struct {
    ID          string
    Name        string
    Description string
    Date        time.Time
    Location    string
    Status      Status
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Domain events (internal, not persisted)
    events []DomainEvent
}

type Status string

const (
    StatusDraft     Status = "draft"
    StatusPublished Status = "published"
    StatusCancelled Status = "cancelled"
)

// Create creates an event and emits a domain event
func (e *Event) Create() {
    e.addEvent(EventCreatedEvent{
        EventID: e.ID,
        Name:    e.Name,
        Date:    e.Date.Format(time.RFC3339),
        // ...
    })
}

// Publish publishes the event and emits a domain event
func (e *Event) Publish() error {
    if e.Status != StatusDraft {
        return errors.New("can only publish draft events")
    }
    
    e.Status = StatusPublished
    e.UpdatedAt = time.Now().UTC()
    
    e.addEvent(EventPublishedEvent{
        EventID: e.ID,
        Name:    e.Name,
        // ...
    })
    
    return nil
}

// Events returns all domain events
func (e Event) Events() []DomainEvent {
    return e.events
}

// ClearEvents clears all domain events
func (e *Event) ClearEvents() {
    e.events = nil
}

// addEvent adds a domain event (internal method)
func (e *Event) addEvent(domainEvent DomainEvent) {
    e.events = append(e.events, domainEvent)
}
```

### 2. `domain_events.go` - Domain Events

```go
package event

// DomainEvent represents a domain event interface
// This is the interface for ALL domain events in this package
type DomainEvent interface {
    Type() string
}

// EventCreatedEvent is a domain event emitted when an Event entity is created
// Note: This is a DOMAIN EVENT, not the Event entity itself
type EventCreatedEvent struct {
    EventID     string
    Name        string
    Description string
    Date        string
    Location    string
    CreatedAt   string
}

func (e EventCreatedEvent) Type() string {
    return "EventCreated"  // Domain event type name
}

// EventPublishedEvent is a domain event emitted when an Event entity is published
type EventPublishedEvent struct {
    EventID string
    Name    string
    Date    string
    PublishedAt string
}

func (e EventPublishedEvent) Type() string {
    return "EventPublished"
}

// EventCancelledEvent is a domain event emitted when an Event entity is cancelled
type EventCancelledEvent struct {
    EventID     string
    Name        string
    CancelledAt string
}

func (e EventCancelledEvent) Type() string {
    return "EventCancelled"
}
```

## Phân biệt rõ ràng

### Event Entity (Business Model)
```go
// domain/event/event.go
type Event struct {
    ID          string  // Business entity
    Name        string
    Date        time.Time
    // ...
}
```

### Domain Events (Business Events)
```go
// domain/event/domain_events.go
type EventCreatedEvent struct {  // Domain event
    EventID string  // Reference to Event entity
    Name    string
    // ...
}
```

## So sánh với Order Domain

### Order Domain (không có conflict):
```
domain/order/
├── order.go    → Order Entity
└── events.go   → Domain Events (OK, không conflict)
```

### Event Domain (có conflict):
```
domain/event/
├── event.go           → Event Entity (business model)
└── domain_events.go   → Domain Events (rõ ràng, tránh nhầm lẫn)
```

## Tại sao dùng `domain_events.go`?

1. **Rõ ràng**: Biết ngay đây là Domain Events, không phải Event Entity
2. **Tránh nhầm lẫn**: Không bị confuse với Event Entity
3. **Consistent**: Pattern nhất quán khi có entity tên "Event"

## Alternative: Nếu muốn ngắn gọn hơn

Có thể dùng `events.go` nếu:
- Code review rõ ràng
- Team đã quen với pattern
- Comment rõ ràng trong code

```
domain/event/
├── event.go    → Event Entity (business model)
└── events.go   → Domain Events (cần comment rõ ràng)
```

**Nhưng khuyến nghị: Dùng `domain_events.go` để rõ ràng nhất.**

