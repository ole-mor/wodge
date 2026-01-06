package monitor

import (
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// EventType defines the type of event
type EventType string

const (
	TypeRequest  EventType = "REQUEST"
	TypePostgres EventType = "POSTGRES"
	TypeRedis    EventType = "REDIS"
	TypeRabbitMQ EventType = "RABBITMQ"
)

// Event represents a monitoring event
type Event struct {
	Timestamp time.Time   `json:"timestamp"`
	Type      EventType   `json:"type"`
	Payload   interface{} `json:"payload"`
}

type Broadcaster struct {
	clients map[chan Event]bool
	mu      sync.Mutex
}

var Bus = &Broadcaster{
	clients: make(map[chan Event]bool),
}

func (b *Broadcaster) Subscribe() chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 100) // buffer to prevent blocking
	b.clients[ch] = true
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[ch]; ok {
		delete(b.clients, ch)
		close(ch)
	}
}

func (b *Broadcaster) Publish(eventType EventType, payload interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	event := Event{
		Timestamp: time.Now(),
		Type:      eventType,
		Payload:   payload,
	}
	for ch := range b.clients {
		select {
		case ch <- event:
		default:
			// Drop event if client is too slow
		}
	}
}

// Handler strictly for the Monitor CLI
func Handler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	clientChan := Bus.Subscribe()
	defer Bus.Unsubscribe(clientChan)

	c.Stream(func(w io.Writer) bool {
		event, ok := <-clientChan
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}

// Middleware to capture HTTP requests
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip monitor events themselves to avoid infinite loop of noise
		if c.Request.URL.Path == "/wodge/monitor/events" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		Bus.Publish(TypeRequest, map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration.String(),
			"ip":       c.ClientIP(),
		})
	}
}
