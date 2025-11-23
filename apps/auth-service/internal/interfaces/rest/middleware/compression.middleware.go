package middleware

import (
	"compress/gzip"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

// CompressionMiddleware creates a compression middleware
// Compresses response if client accepts gzip encoding
func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client accepts gzip
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		// Get gzip writer from pool
		gz := gzipWriterPool.Get().(*gzip.Writer)
		defer func() {
			gz.Reset(nil)
			gzipWriterPool.Put(gz)
		}()

		// Create response writer wrapper
		writer := &gzipResponseWriter{
			ResponseWriter: c.Writer,
			gzipWriter:     gz,
		}
		c.Writer = writer

		// Set headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		c.Next()

		// Flush and close gzip writer
		if writer.gzipWriter != nil {
			writer.gzipWriter.Close()
		}
	}
}

// gzipResponseWriter wraps gin.ResponseWriter to compress output
type gzipResponseWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	// Set content type if not set
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// Reset gzip writer to write to response
	w.gzipWriter.Reset(w.ResponseWriter)
	return w.gzipWriter.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

