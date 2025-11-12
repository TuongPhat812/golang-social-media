# Luồng Tính Năng: Tạo Tin Nhắn Chat

Danh sách dưới đây giúp bạn hiểu toàn bộ pipeline `POST /chat/messages` bằng tiếng Việt.

## 1. Luồng chính (Main Feature)
1. **Định nghĩa route** – `apps/gateway/internal/interfaces/rest/handlers.go`
   - Hàm handler validate payload (`createMessageRequest`) rồi gọi service tầng application.
2. **Service tầng application** – `apps/gateway/internal/application/messages/service.go`
   - Cài đặt interface `Service`. Gọi `ChatClient.CreateMessage` và ánh xạ protobuf response sang domain entity của gateway.
3. **gRPC client** – `apps/gateway/internal/infrastructure/grpc/chat/client.go`
   - Khởi tạo `pkg/gen/chat/v1.ChatServiceClient`, load `CHAT_SERVICE_ADDR` qua `pkg/config`.
4. **Bootstrap server** – `apps/chat-service/cmd/chat-service/main.go`
   - Tải env, mở kết nối GORM tới Postgres (`CHAT_DATABASE_DSN`), khởi tạo Kafka publisher, đăng ký gRPC handler (migrations chạy thủ công bằng `go run ./cmd/migrate`).
5. **gRPC handler** – `apps/chat-service/internal/interfaces/grpc/chat/handler.go`
   - Nhận request (protobuf), gọi service ứng dụng.
6. **Use case / Application service** – `apps/chat-service/internal/application/messages/service.go`
   - Tạo domain `Message`, lưu vào database qua repository, sau đó phát event.

==============================

## 2. Lưu trữ & Side Effect (Kafka & Consumer)
7. **Repository** – `apps/chat-service/internal/infrastructure/persistence/message_repository.go`
   - Triển khai `messages.Repository` bằng GORM.

8. **Kafka publisher** – `apps/chat-service/internal/infrastructure/eventbus/kafka_publisher.go`
   - Triển khai `EventPublisher`, gửi `events.ChatCreated` lên topic `chat.created`.
9. **Notification consumer** – `apps/notification-service/internal/infrastructure/eventbus/subscriber.go`
   - Lắng nghe `chat.created`, forward cho service ứng dụng.
10. **Notification application** – `apps/notification-service/internal/application/notifications/service.go`
   - Tạo domain notification, emit `events.NotificationCreated`.
11. **Notification publisher** – `apps/notification-service/internal/infrastructure/eventbus/kafka_publisher.go`
    - Gửi event mới tới topic `notification.created`.
12. **Socket listeners** – `apps/socket-service/internal/infrastructure/eventbus/listener.go`
    - Có 2 consumer: `chat.created` và `notification.created`.
13. **Socket application service** – `apps/socket-service/internal/application/events/service.go`
    - Log và forward event tới WebSocket hub.
14. **WebSocket hub** – `apps/socket-service/internal/interfaces/socket/hub.go`
    - (Hiện tại) chỉ log broadcast, có thể mở rộng để push dữ liệu tới client.

## 3. Contract & cấu hình dùng chung
15. **Protobuf** – `proto/chat/v1/chat_service.proto` + `pkg/gen/chat/v1`
    - Định nghĩa hợp đồng chéo ngôn ngữ.
16. **Events** – `pkg/events/chat.go`, `pkg/events/notification.go`
    - Payload trung lập dùng cho Kafka.
17. **Environment** – `pkg/config/env.go`
    - Các biến như `KAFKA_BROKERS`, `<SERVICE>_PORT`, `CHAT_DATABASE_DSN`.
18. **Compose** – `docker-compose.infra.yml`, `docker-compose.app.yml`
    - Kiểm tra địa chỉ broker (`kafka:9092` trong container, `localhost:9094` ngoài host) và Postgres (`gsm-postgres:5432` / `localhost:5432`).
