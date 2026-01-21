package services

import (
	"context"
	"io"
)

// DatabaseService defines the interface for database operations (e.g. Postgres)
type DatabaseService interface {
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	// We could add higher level CRUD here later
}

// CacheService defines the interface for cache operations (e.g. Redis)
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Delete(ctx context.Context, key string) error
}

// QueueService defines the interface for message queue operations (e.g. RabbitMQ)
type QueueService interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error
}

// QastService defines the interface for interacting with the QAST API
type QastService interface {
	Ask(ctx context.Context, query, userId, expertise string) (string, []string, error)
	IngestGraph(ctx context.Context, text, userId string) (interface{}, error)
	SecureChat(ctx context.Context, text, userId, token string) (io.ReadCloser, error)
}
