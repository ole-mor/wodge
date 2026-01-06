package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
	"wodge/internal/monitor"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp  string      `json:"timestamp"`
	Method     string      `json:"method"`
	Path       string      `json:"path"`
	Status     int         `json:"status"`
	DurationMs int64       `json:"duration_ms"`
	IP         string      `json:"ip"`
	Body       interface{} `json:"body,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

// RequestLogger returns a middleware that logs detailed request/response info.
// It also taps into the monitor bus.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for the monitor stream itself to prevent loops
		if c.Request.URL.Path == "/wodge/monitor/events" {
			c.Next()
			return
		}

		start := time.Now()

		// Capture Request Body
		var requestBody interface{}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Read body
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// Restore it for next handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Try to parse JSON
			if len(bodyBytes) > 0 {
				_ = json.Unmarshal(bodyBytes, &requestBody)
			}
		}

		// Capture Response Body
		// We need a custom writer to peek at the body
		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		duration := time.Since(start).Milliseconds()

		// Try to parse response body if JSON
		var responseBody interface{}
		_ = json.Unmarshal(w.body.Bytes(), &responseBody)

		entry := LogEntry{
			Timestamp:  start.Format(time.RFC3339),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Status:     c.Writer.Status(),
			DurationMs: duration,
			IP:         c.ClientIP(),
			Body:       requestBody,
			Response:   responseBody,
		}

		// 1. Emit to Monitor Bus (Visualization)
		monitor.Bus.Publish(monitor.TypeRequest, entry)

		// 2. Standard Stdout Log (for container/system logs)
		// We could use slog or zap here, but fmt is fine for now
		// logJSON, _ := json.Marshal(entry)
		// fmt.Println(string(logJSON))
	}
}

// responseBodyWriter is a wrapper to capture the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
