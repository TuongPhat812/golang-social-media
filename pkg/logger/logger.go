package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

const (
	defaultLogDir = "logs"
	logFileName   = "app.log"
)

var (
	once   sync.Once
	global zerolog.Logger
)

func getLogger() *zerolog.Logger {
	once.Do(func() {
		dir := os.Getenv("LOG_OUTPUT_DIR")
		if dir == "" {
			dir = defaultLogDir
		}
		writer := []io.Writer{os.Stdout}
		if fileWriter, err := openFileWriter(dir); err == nil {
			writer = append(writer, fileWriter)
		}
		multi := zerolog.MultiLevelWriter(writer...)
		logger := zerolog.New(multi).With().Timestamp().Logger()
		global = logger
	})
	return &global
}

func openFileWriter(dir string) (io.Writer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, logFileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, getLogger())
}

func FromContext(ctx context.Context) *zerolog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*zerolog.Logger); ok && logger != nil {
		return logger
	}
	return getLogger()
}

func Info() *zerolog.Event {
	return getLogger().Info()
}

func Error() *zerolog.Event {
	return getLogger().Error()
}

func Debug() *zerolog.Event {
	return getLogger().Debug()
}
