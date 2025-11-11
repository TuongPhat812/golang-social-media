module github.com/myself/golang-social-media/apps/chat-service

go 1.25.3

require (
	github.com/myself/golang-social-media/pkg v0.0.0
	github.com/segmentio/kafka-go v0.4.45
	google.golang.org/grpc v1.76.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

replace github.com/myself/golang-social-media/pkg => ../../pkg
