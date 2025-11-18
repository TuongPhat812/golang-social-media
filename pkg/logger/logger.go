package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

const (
	defaultLogDir = "logs"
	logFileName   = "app.log"
	defaultModule = "app"
	timeLayout    = "2006-01-02 15:04:05"
)

var (
	once       sync.Once
	global     zerolog.Logger
	moduleName = os.Getenv("LOG_MODULE")
	moduleMu   sync.RWMutex
)

// SetModule defines the default module name that will be attached to every log line.
// It should be called during service bootstrap before any logging happens.
func SetModule(name string) {
	moduleMu.Lock()
	moduleName = name
	moduleMu.Unlock()
}

// Component returns a logger enriched with the provided component name.
// Example usage: logger.Component("eventbus").Info().Msg("example")
func Component(name string) *zerolog.Logger {
	l := getLogger().With().Str("component", name).Logger()
	return &l
}

func getLogger() *zerolog.Logger {
	once.Do(func() {
		dir := os.Getenv("LOG_OUTPUT_DIR")
		if dir == "" {
			dir = defaultLogDir
		}

		consoleWriter := newConsoleWriter(os.Stdout, !isNoColor())

		writers := []io.Writer{consoleWriter}
		if fileWriter, err := openFileWriter(dir); err == nil {
			writers = append(writers, fileWriter)
		}

		module := currentModule()

		multi := zerolog.MultiLevelWriter(writers...)
		logger := zerolog.New(multi).
			With().
			Timestamp().
			Str("module", module).
			Logger()
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

func currentModule() string {
	moduleMu.RLock()
	defer moduleMu.RUnlock()
	if strings.TrimSpace(moduleName) != "" {
		return moduleName
	}
	return defaultModule
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

type consoleWriter struct {
	out        io.Writer
	timeFormat string
	useColor   bool
}

func newConsoleWriter(out io.Writer, color bool) io.Writer {
	return &consoleWriter{
		out:        out,
		timeFormat: timeLayout,
		useColor:   color,
	}
}

func (cw *consoleWriter) Write(p []byte) (int, error) {
	var evt map[string]interface{}
	if err := json.Unmarshal(p, &evt); err != nil {
		n, err := cw.out.Write(p)
		if err != nil {
			return n, err
		}
		return len(p), nil
	}

	ts := cw.extractString(evt, zerolog.TimestampFieldName)
	if ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			ts = t.Format(cw.timeFormat)
		}
	}

	level := strings.ToUpper(cw.extractString(evt, zerolog.LevelFieldName))
	message := cw.extractString(evt, zerolog.MessageFieldName)
	module := cw.extractString(evt, "module")
	component := cw.extractString(evt, "component")

	line := bytes.Buffer{}

	if ts != "" {
		line.WriteString(cw.style(ts, false))
		line.WriteByte(' ')
	}

	if level != "" {
		line.WriteString(cw.style(fmt.Sprintf("[%s]", level), true, levelColor(level)))
		line.WriteByte(' ')
	}

	if module != "" {
		line.WriteString(cw.style(fmt.Sprintf("[%s]", module), true, colorCyan))
		line.WriteByte(' ')
	}

	if component != "" {
		line.WriteString(cw.style(fmt.Sprintf("[%s]", component), true, colorMagenta))
		line.WriteByte(' ')
	}

	if message != "" {
		line.WriteString(message)
	}

	if len(evt) > 0 {
		keys := make([]string, 0, len(evt))
		for k := range evt {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			line.WriteByte(' ')
			line.WriteString(k)
			line.WriteByte('=')
			line.WriteString(cw.formatValue(evt[k]))
		}
	}

	line.WriteByte('\n')
	n, err := cw.out.Write(line.Bytes())
	if err != nil {
		return n, err
	}
	return len(p), nil
}

func (cw *consoleWriter) extractString(evt map[string]interface{}, key string) string {
	value, ok := evt[key]
	if !ok {
		return ""
	}
	delete(evt, key)

	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case []byte:
		return string(v)
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", v), "0"), ".")
	case nil:
		return ""
	default:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v)
	}
}

func (cw *consoleWriter) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []interface{}:
		return cw.sliceToString(v)
	default:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v)
	}
}

func (cw *consoleWriter) sliceToString(values []interface{}) string {
	strs := make([]string, len(values))
	for i, val := range values {
		strs[i] = cw.formatValue(val)
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

func (cw *consoleWriter) style(text string, bold bool, color ...string) string {
	if text == "" {
		return text
	}

	var builder strings.Builder
	resetNeeded := false

	if bold {
		builder.WriteString(styleBold)
		resetNeeded = true
	}

	if cw.useColor && len(color) > 0 && color[0] != "" {
		builder.WriteString(color[0])
		resetNeeded = true
	}

	builder.WriteString(text)

	if resetNeeded {
		builder.WriteString(colorReset)
	}

	return builder.String()
}

const (
	colorReset   = "\033[0m"
	styleBold    = "\033[1m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

func levelColor(level string) string {
	switch level {
	case "DEBUG":
		return colorBlue
	case "INFO":
		return colorGreen
	case "WARN":
		return colorYellow
	case "ERROR":
		return colorRed
	case "FATAL":
		return colorRed
	case "PANIC":
		return colorRed
	default:
		return colorWhite
	}
}

func isNoColor() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return true
	}
	return false
}
