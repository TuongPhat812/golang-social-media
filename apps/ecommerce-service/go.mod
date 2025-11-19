module golang-social-media/apps/ecommerce-service

go 1.25.3

require (
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.16.0
	github.com/rs/zerolog v1.32.0
	github.com/segmentio/kafka-go v0.4.45
	golang-social-media/pkg v0.0.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.9
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.31.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

replace golang-social-media/pkg => ../../pkg
