# Domain-Driven Design (DDD) Concepts - Implementation Status

## Tổng quan về DDD

Domain-Driven Design là một cách tiếp cận phát triển phần mềm tập trung vào **domain logic** (business logic) và **ubiquitous language** (ngôn ngữ chung giữa dev và domain experts).

## Các khái niệm DDD chính

### 1. **Domain Layer (Tầng Domain)**

#### ✅ **Entities (Đã implement)**
**Định nghĩa:** Objects có identity và lifecycle, được xác định bởi ID.

**Trong repo:**
- `Message` (chat-service) - có ID, SenderID, ReceiverID, Content
- `User` (auth-service, notification-service) - có ID, Email, Name
- `Notification` (notification-service) - có ID, UserID, Type, Title, Body

**Đặc điểm:**
- Có identity (ID)
- Có business logic (Validate(), Create(), MarkAsRead())
- Có domain events
- Immutable trong một số trường hợp

#### ❌ **Value Objects (Chưa implement)**
**Định nghĩa:** Objects không có identity, được xác định bởi giá trị của chúng.

**Ví dụ có thể có:**
- `Email` - validate format, normalize
- `Money` - amount + currency, không thể thay đổi
- `Address` - street, city, country
- `MessageContent` - validate length, sanitize

**Lợi ích:**
- Encapsulate validation logic
- Immutable by design
- Reusable across entities
- Type safety (không thể nhầm Email với string)

**Ví dụ implementation:**
```go
type Email struct {
    value string
}

func NewEmail(s string) (Email, error) {
    // Validate format
    if !isValidEmail(s) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: strings.ToLower(s)}, nil
}

func (e Email) String() string {
    return e.value
}
```

#### ⚠️ **Domain Services (Chưa rõ ràng)**
**Định nghĩa:** Logic không thuộc về một Entity cụ thể, nhưng thuộc về domain.

**Có thể có:**
- `MessageRoutingService` - quyết định route message đến đâu
- `NotificationPriorityService` - tính priority của notification
- `UserMatchingService` - match users dựa trên criteria

**Hiện tại:** Logic này có thể đang nằm trong Commands hoặc Entities.

#### ✅ **Domain Events (Đã implement)**
**Định nghĩa:** Events được emit từ domain layer khi có business events xảy ra.

**Trong repo:**
- `MessageCreatedEvent`
- `UserCreatedEvent`
- `NotificationCreatedEvent`
- `NotificationReadEvent`

**Đặc điểm:**
- Emit từ domain entities
- Immutable
- Represent business events
- Dispatched qua Event Dispatcher

#### ❌ **Aggregates & Aggregate Roots (Chưa implement)**
**Định nghĩa:** 
- **Aggregate:** Cluster of entities và value objects được quản lý như một unit
- **Aggregate Root:** Entity chính của aggregate, là entry point để access aggregate

**Ví dụ có thể có:**
- `ChatAggregate` với root là `Chat` (có nhiều `Message`)
- `UserAggregate` với root là `User` (có `Profile`, `Settings`)
- `NotificationAggregate` với root là `Notification` (có `Metadata`, `Actions`)

**Lợi ích:**
- Consistency boundary
- Transaction boundary
- Encapsulation
- Clear ownership

**Hiện tại:** Mỗi entity đang được quản lý độc lập, chưa có aggregate boundaries rõ ràng.

---

### 2. **Application Layer (Tầng Application)**

#### ✅ **Commands (Đã implement)**
**Định nghĩa:** Operations thay đổi state của system.

**Trong repo:**
- `CreateMessageCommand`
- `RegisterUserCommand`
- `CreateNotificationCommand`
- `MarkNotificationReadCommand`

**Đặc điểm:**
- Có Execute() method
- Có contracts/interfaces
- Return domain entities hoặc DTOs
- Side effects (thay đổi state)

#### ✅ **Queries (Đã implement)**
**Định nghĩa:** Operations đọc data, không thay đổi state.

**Trong repo:**
- `GetUserProfileQuery`
- `GetNotificationsQuery`

**Đặc điểm:**
- Có Execute() method
- Có contracts/interfaces
- Read-only
- Return DTOs hoặc domain entities

#### ✅ **Event Handlers (Đã implement)**
**Định nghĩa:** Handle domain events và transform sang infrastructure events.

**Trong repo:**
- `MessageCreatedHandler`
- `UserCreatedHandler`
- `NotificationCreatedHandler`
- `NotificationReadHandler`

**Đặc điểm:**
- Transform domain events → infrastructure events
- Abstraction over infrastructure (EventBrokerPublisher)
- Không có business logic

#### ✅ **Event Dispatcher (Đã implement)**
**Định nghĩa:** Dispatch domain events đến registered handlers.

**Trong repo:**
- `event_dispatcher.Dispatcher`
- Register handlers by event type
- Dispatch events asynchronously

#### ❌ **Application Services (Chưa rõ ràng)**
**Định nghĩa:** Orchestrate multiple domain operations, coordinate between aggregates.

**Có thể có:**
- `ChatOrchestrationService` - coordinate tạo message + send notification
- `UserOnboardingService` - coordinate register + send welcome email + create profile

**Hiện tại:** Logic này có thể đang nằm trong Commands.

#### ❌ **Specification Pattern (Chưa implement)**
**Định nghĩa:** Encapsulate business rules dưới dạng reusable specifications.

**Ví dụ có thể có:**
- `IsAdultUserSpec` - user phải >= 18 tuổi
- `CanSendMessageSpec` - check user có thể gửi message không
- `IsValidNotificationSpec` - validate notification rules

**Lợi ích:**
- Reusable business rules
- Composable (AND, OR, NOT)
- Testable
- Clear business intent

---

### 3. **Infrastructure Layer (Tầng Infrastructure)**

#### ✅ **Repositories (Đã implement)**
**Định nghĩa:** Abstraction over data persistence.

**Trong repo:**
- `MessageRepository` (PostgreSQL)
- `NotificationRepository` (ScyllaDB)
- `UserRepository` (Memory)

**Đặc điểm:**
- Implement domain repository interfaces
- Handle persistence details
- Transform domain ↔ persistence models

#### ⚠️ **Repository Pattern (Chưa hoàn chỉnh)**
**Vấn đề hiện tại:**
- Repository interfaces nằm trong application layer (messages.Repository)
- Nên nằm trong domain layer hoặc application layer với contracts rõ ràng hơn

**Nên có:**
- Domain repository interfaces (trong domain hoặc application/contracts)
- Infrastructure implementations (trong infrastructure)

#### ✅ **Event Bus (Đã implement)**
**Định nghĩa:** Infrastructure cho event publishing/subscribing.

**Trong repo:**
- Kafka Publisher (với contracts)
- Kafka Subscriber (với contracts)
- EventBrokerAdapter (abstraction layer)

#### ❌ **Outbox Pattern (Chưa implement)**
**Định nghĩa:** Pattern đảm bảo events được publish reliably.

**Cách hoạt động:**
1. Save events vào outbox table trong cùng transaction với business data
2. Background worker đọc từ outbox và publish
3. Mark as published sau khi thành công

**Lợi ích:**
- Guaranteed delivery
- Transactional consistency
- Retry mechanism

**Hiện tại:** Events được publish trực tiếp, có thể mất nếu service crash.

---

### 4. **Interfaces Layer (Tầng Interfaces)**

#### ✅ **gRPC Handlers (Đã implement)**
**Trong repo:**
- `chat.Handler` - gRPC handlers cho chat-service

#### ✅ **HTTP Handlers (Đã implement)**
**Trong repo:**
- REST handlers trong gateway
- REST handlers trong auth-service

#### ✅ **WebSocket Handlers (Đã implement)**
**Trong repo:**
- `socket.Hub` - WebSocket hub cho socket-service

---

## Các pattern DDD khác

### ❌ **Factory Pattern (Chưa implement)**
**Định nghĩa:** Tạo complex objects hoặc aggregates.

**Ví dụ có thể có:**
- `MessageFactory` - tạo message với validation và domain events
- `NotificationFactory` - tạo notification dựa trên type
- `UserFactory` - tạo user với default settings

**Lợi ích:**
- Encapsulate complex creation logic
- Ensure invariants
- Reusable creation patterns

### ❌ **Saga Pattern (Chưa implement)**
**Định nghĩa:** Coordinate distributed transactions across multiple services.

**Ví dụ có thể có:**
- `CreateUserSaga` - coordinate: create user → send welcome email → create profile
- `SendMessageSaga` - coordinate: create message → send notification → update unread count

**Lợi ích:**
- Distributed transaction management
- Compensation logic
- Eventual consistency

**Hiện tại:** Các operations đang được coordinate qua events, nhưng chưa có explicit saga pattern.

### ❌ **Event Sourcing (Chưa implement)**
**Định nghĩa:** Store events thay vì current state, reconstruct state từ events.

**Cách hoạt động:**
- Store tất cả events (MessageCreated, MessageUpdated, MessageDeleted)
- Reconstruct current state bằng cách replay events
- Event store là source of truth

**Lợi ích:**
- Complete audit trail
- Time travel (xem state tại bất kỳ thời điểm nào)
- Event replay for debugging

**Hiện tại:** Chỉ store current state, không có event store.

### ⚠️ **Bounded Contexts (Có nhưng chưa rõ ràng)**
**Định nghĩa:** Explicit boundaries giữa các domain models.

**Trong repo:**
- `auth-service` - User context (authentication, authorization)
- `chat-service` - Message context
- `notification-service` - Notification context
- `gateway` - Orchestration context

**Đặc điểm:**
- Mỗi service có domain model riêng
- User trong auth-service khác với User trong notification-service
- Communication qua events hoặc APIs

**Có thể cải thiện:**
- Explicit context mapping
- Anti-corruption layers
- Shared kernel (nếu cần)

### ❌ **Anti-Corruption Layer (Chưa implement)**
**Định nghĩa:** Layer bảo vệ domain model khỏi external systems.

**Ví dụ có thể có:**
- Transform external API responses → domain models
- Isolate domain từ third-party services
- Handle versioning và compatibility

**Hiện tại:** Gateway đang call trực tiếp external services, chưa có ACL.

---

## Tóm tắt Implementation Status

### ✅ Đã implement:
1. **Entities** - Domain entities với business logic
2. **Domain Events** - Events từ domain layer
3. **Commands & Queries** - CQRS pattern
4. **Event Handlers** - Transform domain events
5. **Event Dispatcher** - Dispatch events
6. **Repositories** - Data persistence abstraction
7. **Event Bus** - Kafka publisher/subscriber với contracts
8. **Layered Architecture** - Domain, Application, Infrastructure, Interfaces
9. **Contracts/Interfaces** - Dependency inversion
10. **Bootstrap Pattern** - Dependency setup

### ⚠️ Chưa rõ ràng / Cần cải thiện:
1. **Value Objects** - Chưa có explicit value objects
2. **Domain Services** - Logic có thể đang nằm trong Commands
3. **Aggregates** - Chưa có aggregate boundaries rõ ràng
4. **Repository Interfaces** - Nên nằm trong domain/application contracts
5. **Bounded Contexts** - Có nhưng chưa explicit mapping

### ❌ Chưa implement:
1. **Value Objects** - Email, Money, Address, etc.
2. **Aggregates & Aggregate Roots** - ChatAggregate, UserAggregate
3. **Domain Services** - Business logic không thuộc entity
4. **Specification Pattern** - Reusable business rules
5. **Factory Pattern** - Complex object creation
6. **Outbox Pattern** - Reliable event publishing
7. **Saga Pattern** - Distributed transaction coordination
8. **Event Sourcing** - Event store và replay
9. **Anti-Corruption Layer** - Protection từ external systems
10. **Shared Kernel** - Shared domain concepts (nếu cần)

---

## Khuyến nghị cải thiện

### Priority 1 (Quan trọng):
1. **Value Objects** - Email, MessageContent để type safety và validation
2. **Outbox Pattern** - Đảm bảo events không bị mất
3. **Aggregates** - Rõ ràng consistency boundaries

### Priority 2 (Nên có):
4. **Domain Services** - Extract logic không thuộc entity
5. **Specification Pattern** - Reusable business rules
6. **Factory Pattern** - Complex creation logic

### Priority 3 (Nice to have):
7. **Saga Pattern** - Nếu cần distributed transactions
8. **Event Sourcing** - Nếu cần audit trail và time travel
9. **Anti-Corruption Layer** - Nếu có nhiều external integrations

---

## Kết luận

Repo hiện tại đã implement **khoảng 60-70%** các khái niệm DDD cơ bản:
- ✅ Layered architecture rõ ràng
- ✅ Domain entities với business logic
- ✅ Domain events pattern
- ✅ CQRS pattern
- ✅ Repository pattern (cần cải thiện)
- ✅ Event-driven architecture

**Còn thiếu:**
- Value Objects (quan trọng cho type safety)
- Aggregates (quan trọng cho consistency)
- Outbox Pattern (quan trọng cho reliability)
- Một số patterns nâng cao (Saga, Event Sourcing)

Repo đang ở mức **good DDD implementation** với room for improvement.

