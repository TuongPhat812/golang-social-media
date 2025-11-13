## Gateway CQRS Architecture

```
apps/gateway/
├── cmd/gateway/
│   └── main.go                     # composition root
├── internal/
│   ├── application/
│   │   ├── command/
│   │   │   ├── contracts/
│   │   │   │   ├── register_user.command.contract.go
│   │   │   │   ├── login_user.command.contract.go
│   │   │   │   ├── create_message.command.contract.go
│   │   │   │   ├── send_friend_request.command.contract.go
│   │   │   │   ├── accept_friend_request.command.contract.go
│   │   │   │   └── ... (mỗi command một contract)
│   │   │   ├── register_user.command.go
│   │   │   ├── login_user.command.go
│   │   │   ├── create_message.command.go
│   │   │   ├── send_friend_request.command.go
│   │   │   ├── accept_friend_request.command.go
│   │   │   └── dto/
│   │   │       ├── register_user.command.dto.go
│   │   │       ├── login_user.command.dto.go
│   │   │       ├── create_message.command.dto.go
│   │   │       └── friend_request.command.dto.go
│   │   ├── query/
│   │   │   ├── contracts/
│   │   │   │   ├── get_user_profile.query.contract.go
│   │   │   │   ├── get_friend_list.query.contract.go
│   │   │   │   ├── list_chat_threads.query.contract.go
│   │   │   │   └── ... (mỗi query một contract)
│   │   │   ├── get_user_profile.query.go
│   │   │   ├── get_friend_list.query.go
│   │   │   ├── list_chat_threads.query.go
│   │   │   └── dto/
│   │   │       ├── user_profile.query.dto.go
│   │   │       ├── friend_list.query.dto.go
│   │   │       └── chat_thread.query.dto.go
│   │   └── registry/
│   │       └── wiring.go                    # helper to instantiate handlers
│   ├── domain/
│   │   ├── message/
│   │   │   └── entity.go
│   │   └── user/
│   │       └── entity.go
│   ├── infrastructure/
│   │   ├── grpc/
│   │   │   └── chat/
│   │   │       └── client.go
│   │   ├── persistence/
│   │   │   └── read_models/...              # each read model per file
│   │   └── http/
│   │       ├── router.go                    # build gin.Engine using interfaces
│   │       └── middleware/...
│   └── interfaces/
│       └── rest/
│           ├── command/
│           │   ├── register_user.http.handler.go
│           │   ├── login_user.http.handler.go
│           │   ├── create_message.http.handler.go
│           │   ├── send_friend_request.http.handler.go
│           │   ├── accept_friend_request.http.handler.go
│           │   └── contracts/
│           │       ├── register_user.http.contract.go
│           │       ├── login_user.http.contract.go
│           │       ├── create_message.http.contract.go
│           │       ├── send_friend_request.http.contract.go
│           │       └── accept_friend_request.http.contract.go
│           ├── query/
│           │   ├── get_user_profile.http.handler.go
│           │   ├── get_friend_list.http.handler.go
│           │   ├── list_chat_threads.http.handler.go
│           │   └── contracts/
│           │       ├── get_user_profile.http.contract.go
│           │       ├── get_friend_list.http.contract.go
│           │       └── list_chat_threads.http.contract.go
│           └── dto/
│               ├── register_user.http.request.go
│               ├── login_user.http.request.go
│               ├── create_message.http.request.go
│               ├── send_friend_request.http.request.go
│               ├── accept_friend_request.http.request.go
│               ├── user_profile.http.response.go
│               ├── friend_list.http.response.go
│               └── chat_threads.http.response.go
└── pkg/
    └── ...
```

*Each use-case/feature lives in its own file. Interfaces are separated from implementations to keep dependency inversion explicit.* 

