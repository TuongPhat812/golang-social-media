package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"golang-social-media/pkg/logger"
)

var loadOnce sync.Once

// LoadEnv loads environment variables from the provided files (defaults to .env).
func LoadEnv(filenames ...string) {
	loadOnce.Do(func() {
		files := filenames
		if len(files) == 0 {
			files = discoverEnvPaths(".env")
		}

		for _, file := range files {
			if file == "" {
				continue
			}
			if _, err := os.Stat(file); err == nil {
				if err := godotenv.Load(file); err != nil {
					logger.Component("config").
						Error().
						Err(err).
						Str("file", file).
						Msg("config unable to load env file")
				}
				return
			}
		}
	})
}

// GetEnv returns the value of the environment variable or a fallback when unset.
func GetEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// GetEnvInt parses an integer environment variable, returning fallback on error.
func GetEnvInt(key string, fallback int) int {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		logger.Component("config").
			Error().
			Str("key", key).
			Str("value", value).
			Msg("config invalid integer env value")
	}
	return fallback
}

// GetEnvStringSlice splits a comma-separated environment variable into a slice.
func GetEnvStringSlice(key string, fallback []string) []string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		parts := strings.Split(value, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				out = append(out, trimmed)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return fallback
}

func discoverEnvPaths(filename string) []string {
	wd, err := os.Getwd()
	if err != nil {
		return []string{filename}
	}

	var paths []string
	for dir := wd; ; dir = filepath.Dir(dir) {
		paths = append(paths, filepath.Join(dir, filename))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return paths
}
