package main

import (
	"log"
	"net/http"
	"os"
	"wodge/internal/drivers/postgres"
	"wodge/internal/drivers/rabbitmq"
	"wodge/internal/drivers/redis"
	"wodge/internal/monitor"
	"wodge/internal/services"

	"github.com/gin-gonic/gin"
)

// Global services
var (
	db    services.DatabaseService
	cache services.CacheService
	queue services.QueueService
)

func main() {
	// Initialize Services
	initServices()

	r := gin.Default()

	// Add Monitor Middleware
	r.Use(monitor.Middleware())

	// Enable CORS for dev server
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register API endpoints
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "wow"})
	})

	// Monitor Event Stream
	r.GET("/wodge/monitor/events", monitor.Handler)

	// Service Routes
	api := r.Group("/api")
	{
		// Postgres Routes
		api.POST("/postgres/query", handlePostgresQuery)
		api.POST("/postgres/execute", handlePostgresExecute)

		// Redis Routes
		api.GET("/redis/:key", handleRedisGet)
		api.POST("/redis", handleRedisSet)
		api.DELETE("/redis/:key", handleRedisDelete)

		// RabbitMQ Routes
		// Note: Subscribe is streaming/push, simpler to just allow Publish via HTTP for now
		api.POST("/queue/publish", handleQueuePublish)
	}

	log.Println("Starting API server on :8080")
	log.Println("Frontend will access APIs via: http://localhost:5173/api")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func initServices() {
	// Postgres
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn != "" {
		var err error
		db, err = postgres.NewPostgresDriver(dsn)
		if err != nil {
			log.Printf("Failed to init Postgres: %v", err)
		} else {
			log.Println("Postgres connected")
		}
	}

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		pass := os.Getenv("REDIS_PASSWORD")
		var err error
		cache, err = redis.NewRedisDriver(redisAddr, pass, 0)
		if err != nil {
			log.Printf("Failed to init Redis: %v", err)
		} else {
			log.Println("Redis connected")
		}
	}

	// RabbitMQ
	amqpUrl := os.Getenv("RABBITMQ_URL")
	if amqpUrl != "" {
		var err error
		queue, err = rabbitmq.NewRabbitMQDriver(amqpUrl)
		if err != nil {
			log.Printf("Failed to init RabbitMQ: %v", err)
		} else {
			log.Println("RabbitMQ connected")
		}
	}
}

// -- Handlers --

// POST /api/postgres/query { "query": "SELECT...", "args": [...] }
func handlePostgresQuery(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Postgres not configured"})
		return
	}
	var req struct {
		Query string        `json:"query"`
		Args  []interface{} `json:"args"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results, err := db.Query(c.Request.Context(), req.Query, req.Args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// POST /api/postgres/execute { "query": "INSERT...", "args": [...] }
func handlePostgresExecute(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Postgres not configured"})
		return
	}
	var req struct {
		Query string        `json:"query"`
		Args  []interface{} `json:"args"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rows, err := db.Execute(c.Request.Context(), req.Query, req.Args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows_affected": rows})
}

// GET /api/redis/:key
func handleRedisGet(c *gin.Context) {
	if cache == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Redis not configured"})
		return
	}
	key := c.Param("key")
	val, err := cache.Get(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"}) // approximate
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": val})
}

// POST /api/redis { "key": "...", "value": "...", "ttl": 60 }
func handleRedisSet(c *gin.Context) {
	if cache == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Redis not configured"})
		return
	}
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cache.Set(c.Request.Context(), req.Key, req.Value, req.TTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DELETE /api/redis/:key
func handleRedisDelete(c *gin.Context) {
	if cache == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Redis not configured"})
		return
	}
	key := c.Param("key")
	if err := cache.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// POST /api/queue/publish { "topic": "...", "message": "..." }
func handleQueuePublish(c *gin.Context) {
	if queue == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "RabbitMQ not configured"})
		return
	}
	var req struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := queue.Publish(c.Request.Context(), req.Topic, []byte(req.Message)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
