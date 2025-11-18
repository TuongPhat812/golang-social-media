# Domain Events Pattern với CQRS

## Pattern Flow

```
Application Layer (Command)
    ↓
acceptFriendRequestUseCase
    ↓ input DTO
    ↓
create domain model FriendRelationship
    ↓
FriendRelationship.accept() [Domain Logic]
    ↓
FriendRelationship.addEvent(FriendRelationshipCreated) [Domain Event]
    ↓
repo.persist(FriendRelationship) [Persistence]
    ↓
dispatch domain events [Event Dispatching]
    ↓
for event in FriendRelationship.events():
    if event is FriendRelationshipCreated:
        FriendRelationshipCreatedEventHandler.handle(event)
            ↓
Application Layer (Event Handler)
            ↓
Infrastructure Layer
            ↓
kafkaProducer.publish(event)
```

## Key Principles

### 1. **Domain Entity tracks events**
- Domain entity có method `addEvent()` để collect events
- Domain entity có method `events()` để retrieve events
- Events được tạo từ business logic trong domain

### 2. **Events dispatched AFTER persistence**
- Dispatch events sau khi transaction commit thành công
- Đảm bảo data consistency trước khi publish events

### 3. **Event Handlers ở Application Layer**
- Event handlers là application services
- Không phải domain logic, mà là orchestration

### 4. **Infrastructure publishes events**
- Kafka publisher ở infrastructure layer
- Domain không biết về Kafka, chỉ biết về events

## Implementation Example

### Domain Entity với Events

```go
// domain/friend_relationship/entity.go
type FriendRelationship struct {
    ID        string
    UserID    string
    FriendID  string
    Status    Status
    CreatedAt time.Time
    
    // Domain events (internal, not persisted)
    events []DomainEvent
}

func (f *FriendRelationship) accept() error {
    // Business logic
    if f.Status != StatusPending {
        return errors.New("can only accept pending requests")
    }
    
    f.Status = StatusAccepted
    
    // Add domain event
    f.addEvent(FriendRelationshipCreated{
        RelationshipID: f.ID,
        UserID:         f.UserID,
        FriendID:       f.FriendID,
        Status:         f.Status,
    })
    
    return nil
}

func (f *FriendRelationship) addEvent(event DomainEvent) {
    f.events = append(f.events, event)
}

func (f *FriendRelationship) Events() []DomainEvent {
    return f.events
}

func (f *FriendRelationship) ClearEvents() {
    f.events = nil
}
```

### Application Command

```go
// application/command/accept_friend_request.command.go
func (c *acceptFriendRequestCommand) Execute(ctx context.Context, req dto.AcceptFriendRequestRequest) error {
    // 1. Load domain entity
    relationship, err := c.repo.FindByID(ctx, req.RelationshipID)
    if err != nil {
        return err
    }
    
    // 2. Execute domain logic (this adds events internally)
    if err := relationship.accept(); err != nil {
        return err
    }
    
    // 3. Validate
    if err := relationship.Validate(); err != nil {
        return err
    }
    
    // 4. Persist
    if err := c.repo.Save(ctx, relationship); err != nil {
        return err
    }
    
    // 5. Dispatch domain events (AFTER successful persistence)
    events := relationship.Events()
    relationship.ClearEvents() // Clear after dispatch
    
    for _, event := range events {
        if err := c.eventDispatcher.Dispatch(ctx, event); err != nil {
            // Log error but don't fail the command
            // Events can be retried via outbox pattern
            c.log.Error().Err(err).Msg("failed to dispatch domain event")
        }
    }
    
    return nil
}
```

### Event Handler

```go
// application/event_handler/friend_relationship_created.handler.go
func (h *FriendRelationshipCreatedHandler) Handle(ctx context.Context, event FriendRelationshipCreated) error {
    // Transform domain event to infrastructure event
    kafkaEvent := events.FriendRelationshipCreated{
        RelationshipID: event.RelationshipID,
        UserID:         event.UserID,
        FriendID:       event.FriendID,
        Status:         string(event.Status),
        CreatedAt:      time.Now(),
    }
    
    // Publish via infrastructure
    return h.publisher.PublishFriendRelationshipCreated(ctx, kafkaEvent)
}
```

### Event Dispatcher

```go
// application/event_dispatcher/dispatcher.go
type EventDispatcher interface {
    Dispatch(ctx context.Context, event DomainEvent) error
}

type dispatcher struct {
    handlers map[string][]EventHandler
}

func (d *dispatcher) Dispatch(ctx context.Context, event DomainEvent) error {
    eventType := event.Type()
    handlers := d.handlers[eventType]
    
    for _, handler := range handlers {
        if err := handler.Handle(ctx, event); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Benefits

1. **Separation of Concerns**: Domain logic tách biệt với infrastructure
2. **Testability**: Domain logic có thể test độc lập
3. **Flexibility**: Có thể thêm nhiều event handlers mà không thay đổi domain
4. **Consistency**: Events chỉ được dispatch sau khi data đã persist
5. **Decoupling**: Domain không phụ thuộc vào infrastructure

## Transaction Boundaries

⚠️ **Important**: Events nên được dispatch sau khi transaction commit thành công. Nếu cần đảm bảo exactly-once delivery, có thể dùng:
- **Outbox Pattern**: Store events in DB, then publish via separate process
- **Transaction Log Tailing**: Read from DB transaction log
- **Saga Pattern**: For distributed transactions

