module golang-social-media/apps/ecommerce-service

go 1.25.3

require (
	github.com/google/uuid v1.6.0
	github.com/rs/zerolog v1.32.0
	github.com/segmentio/kafka-go v0.4.45
	golang-social-media/pkg v0.0.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.7
)

require (
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
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)

replace golang-social-media/pkg => ../../pkg
