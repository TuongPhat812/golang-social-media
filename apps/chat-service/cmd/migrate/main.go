package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang-social-media/pkg/logger"
)

const (
	migrationsDir = "migrations"
	defaultDSN    = "postgres://chat_user:chat_password@localhost:5432/chat_service?sslmode=disable"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/migrate create <name>")
	fmt.Println("  go run ./cmd/migrate up")
	fmt.Println("  go run ./cmd/migrate down")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create":
		if len(os.Args) < 3 {
			logger.Error().Msg("missing migration name. Example: go run ./cmd/migrate create add_messages_index")
			os.Exit(1)
		}
		name := os.Args[2]
		if err := createMigration(name); err != nil {
			logger.Error().Err(err).Msg("failed to create migration")
			os.Exit(1)
		}
		fmt.Printf("Created migration %s\n", name)
	case "up":
		if err := runMigrations(1); err != nil {
			logger.Error().Err(err).Msg("failed to run migrations up")
			os.Exit(1)
		}
		fmt.Println("Migrations applied")
	case "down":
		if err := runMigrations(-1); err != nil {
			logger.Error().Err(err).Msg("failed to run migrations down")
			os.Exit(1)
		}
		fmt.Println("Migrations rolled back one step")
	default:
		usage()
		os.Exit(1)
	}
}

func createMigration(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("migration name cannot be empty")
	}
	if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format("20060102150405")
	base := fmt.Sprintf("%s_%s", timestamp, name)
	upPath := filepath.Join(migrationsDir, base+".up.sql")
	downPath := filepath.Join(migrationsDir, base+".down.sql")

	if err := os.WriteFile(upPath, []byte("-- write your UP migration here\n"), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(downPath, []byte("-- write your DOWN migration here\n"), 0o644); err != nil {
		return err
	}
	return nil
}

func runMigrations(direction int) error {
	dsn := os.Getenv("CHAT_DATABASE_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	m, err := migrate.New("file://"+absDir, dsn)
	if err != nil {
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Error().Err(srcErr).Msg("migration source close error")
		}
		if dbErr != nil {
			logger.Error().Err(dbErr).Msg("migration database close error")
		}
	}()

	switch {
	case direction > 0:
		err = m.Up()
	case direction < 0:
		err = m.Steps(-1)
	default:
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info().Msg("migration: no changes to apply")
		return nil
	}
	return err
}
