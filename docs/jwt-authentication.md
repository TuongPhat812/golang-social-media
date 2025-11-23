# JWT Authentication Implementation

## Overview

Đã implement JWT authentication cho endpoint tạo chat message. User phải có JWT token hợp lệ để tạo message.

## Changes

### 1. Auth Service

#### JWT Service (`apps/auth-service/internal/infrastructure/jwt/jwt.service.go`)
- `GenerateToken(userID string)`: Generate JWT token với user_id claim
- `ValidateToken(tokenString string)`: Validate và extract user_id từ token
- Configurable expiration time (default: 24 hours)

#### Login Command Update
- Thay thế simple token bằng JWT token
- Sử dụng JWT service để generate token
- Token chứa user_id trong claims

#### Bootstrap Update
- Setup JWT service với secret và expiration từ env
- Environment variables:
  - `JWT_SECRET`: Secret key để sign tokens (default: "your-secret-key-change-in-production")
  - `JWT_EXPIRATION_HOURS`: Token expiration in hours (default: 24)

### 2. Gateway Service

#### JWT Middleware (`apps/gateway/internal/infrastructure/middleware/jwt.middleware.go`)
- Validate JWT token từ Authorization header
- Format: `Authorization: Bearer <token>`
- Extract user_id từ token claims
- Set user_id vào Gin context để handlers sử dụng

#### Router Update
- Apply JWT middleware cho protected routes
- `/chat/messages` endpoint giờ yêu cầu JWT authentication
- Public routes (`/auth/*`) không cần JWT

#### Create Message Handler Update
- Lấy `sender_id` từ JWT context (không cần trong request body nữa)
- Request body chỉ cần `receiverId` và `content`
- `senderId` được tự động lấy từ authenticated user

## API Changes

### Before:
```json
POST /chat/messages
{
  "senderId": "user-1",
  "receiverId": "user-2",
  "content": "Hello"
}
```

### After:
```json
POST /chat/messages
Authorization: Bearer <jwt_token>
{
  "receiverId": "user-2",
  "content": "Hello"
}
```

## Environment Variables

### Auth Service:
```bash
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24
```

### Gateway:
```bash
JWT_SECRET=your-secret-key-change-in-production  # Must match auth service
```

## Flow

1. User login → Auth service generate JWT token
2. Client gửi request với JWT token trong Authorization header
3. Gateway JWT middleware validate token
4. Extract user_id từ token → Set vào context
5. Create message handler lấy sender_id từ context
6. Create message với authenticated user as sender

## Security Benefits

- ✅ User không thể fake sender_id (lấy từ JWT)
- ✅ Token có expiration (default 24h)
- ✅ Token được sign với secret key
- ✅ Stateless authentication (không cần token store)

## Next Steps

1. Add token refresh mechanism (optional)
2. Add role-based access control (optional)
3. Add token blacklist for logout (optional)
4. Use stronger secret key in production

