package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wodge/internal/drivers/postgres"
	"wodge/internal/drivers/qast"
	"wodge/internal/drivers/rabbitmq"
	"wodge/internal/drivers/redis"
	"wodge/internal/middleware"
	"wodge/internal/monitor"
	"wodge/internal/services"

	"github.com/gin-gonic/gin"
)

// Global services
var (
	db      services.DatabaseService
	cache   services.CacheService
	queue   services.QueueService
	qastSvc services.QastService
)

// Start starts the Wodge API server
func Start(port int) {
	// Print debug info about env vars
	log.Printf("DEBUG: POSTGRES_DSN=%s", os.Getenv("POSTGRES_DSN"))
	log.Printf("DEBUG: REDIS_ADDR=%s", os.Getenv("REDIS_ADDR"))

	// Initialize Services
	initServices()

	r := gin.Default()

	// Add Request Logging Middleware
	r.Use(middleware.RequestLogger())

	// Enable CORS for dev server
	r.Use(func(c *gin.Context) {
		// Dynamic ORIGIN support for development
		// In production, this should be stricter, but for dev tool we can trust localhost
		origin := c.Request.Header.Get("Origin")
		// Check if origin is localhost or 127.0.0.1
		// Simplest for now: Allow all localhost ports
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Fallback
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		}

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
		c.JSON(200, gin.H{"status": "ok"})
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

		// QAST Routes
		api.POST("/qast/ask", handleQastAsk)
		api.POST("/qast/ingest", handleQastIngest)
	}

	log.Printf("Starting Wodge API server on :%d\n", port)
	log.Println("Frontend will access APIs via: http://localhost:5173/api")

	// Format address
	addr := fmt.Sprintf(":%d", port)

	if err := r.Run(addr); err != nil {
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
			log.Printf("ERROR: Failed to init Postgres: %v", err)
			db = nil // Ensure strictly nil
		} else {
			log.Println("Postgres connected successfully")
		}
	} else {
		log.Println("POSTGRES_DSN is empty, skipping Postgres init")
	}

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		pass := os.Getenv("REDIS_PASSWORD")
		var err error
		cache, err = redis.NewRedisDriver(redisAddr, pass, 0)
		if err != nil {
			log.Printf("ERROR: Failed to init Redis: %v", err)
			cache = nil
		} else {
			log.Println("Redis connected successfully")
		}
	} else {
		log.Println("REDIS_ADDR is empty, skipping Redis init")
	}

	// RabbitMQ
	amqpUrl := os.Getenv("RABBITMQ_URL")
	if amqpUrl != "" {
		var err error
		queue, err = rabbitmq.NewRabbitMQDriver(amqpUrl)
		if err != nil {
			log.Printf("ERROR: Failed to init RabbitMQ: %v", err)
			queue = nil
		} else {
			log.Println("RabbitMQ connected successfully")
		}
	}

	// QAST
	qastURL := os.Getenv("QAST_URL")
	if qastURL != "" {
		qastSvc = qast.NewQastDriver(qastURL)
		log.Println("QAST driver initialized")
	} else {
		log.Println("QAST_URL is empty, skipping QAST init")
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

// POST /api/qast/ask { "query": "..." }
func handleQastAsk(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	var req struct {
		Query          string `json:"query"`
		UserID         string `json:"user_id"`
		ExpertiseLevel string `json:"expertise_level"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer, context, err := qastSvc.Ask(c.Request.Context(), req.Query, req.UserID, req.ExpertiseLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": answer, "context": context})
}

// POST /api/qast/ingest { "text": "..." }
func handleQastIngest(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	var req struct {
		Text   string `json:"text"`
		UserID string `json:"user_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := qastSvc.IngestGraph(c.Request.Context(), req.Text, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": result})
}
