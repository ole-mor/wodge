package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"wodge/internal/drivers/astauth"
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
	db         services.DatabaseService
	cache      services.CacheService
	queue      services.QueueService
	qastSvc    services.QastService
	astAuthSvc *astauth.AstAuthDriver
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
		api.POST("/qast/ingest/async", handleQastIngestAsync)
		api.POST("/qast/chat", handleQastSecureChat)

		// History Routes (Qast Proxy)
		api.POST("/history/sessions", handleHistoryCreateSession)
		api.GET("/history/sessions", handleHistoryGetSessions)
		api.GET("/history/sessions/:id", handleHistoryGetSession)
		api.DELETE("/history/sessions/:id", handleHistoryDeleteSession)

		// Auth Routes
		api.POST("/auth/login", handleAuthLogin)
		api.POST("/auth/register", handleAuthRegister)
		api.POST("/auth/refresh", handleAuthRefresh)
		api.POST("/auth/verify", handleAuthVerify)
		api.GET("/users/me", handleAuthVerify) // Alias for verify
		api.POST("/auth/logout", handleAuthLogout)
		api.GET("/users/search", handleUsersSearch)

		// Share Route
		api.POST("/history/sessions/:id/share", handleHistoryShareSession)

		// Context Routes
		api.PUT("/context/:id", handleContextUpdate)
		api.GET("/context/:id", handleContextGet)
	}

	log.Printf("Starting Wodge API server on :%d\n", port)
	log.Println("--- WODGE SERVER VERSION: CONTEXT_UPDATE_PATCHED ---")
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
		apiKey := os.Getenv("QAST_API_KEY")
		qastSvc = qast.NewQastDriver(qastURL, apiKey)
		log.Println("QAST driver initialized")
	} else {
		log.Println("QAST_URL is empty, skipping QAST init")
	}

	// AstAuth
	astAuthURL := os.Getenv("ASTAUTH_URL")
	if astAuthURL != "" {
		astAuthSvc = astauth.NewAstAuthDriver(astAuthURL)
		log.Println("AstAuth driver initialized")
	} else {
		log.Println("ASTAUTH_URL is empty, skipping AstAuth init")
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

// POST /api/qast/ingest/async { "text": "..." }
func handleQastIngestAsync(c *gin.Context) {
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

	// Run in background
	go func() {
		// Create a background context since request context will be cancelled
		ctx := context.Background()
		log.Printf("Starting async ingest for user %s...", req.UserID)
		_, err := qastSvc.IngestGraph(ctx, req.Text, req.UserID)
		if err != nil {
			log.Printf("Async ingest failed: %v", err)
		} else {
			log.Printf("Async ingest completed for user %s", req.UserID)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted", "message": "Ingestion started in background"})
}

func handleQastSecureChat(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	var req struct {
		Text            string `json:"text"`
		UserID          string `json:"user_id"`
		SessionID       string `json:"session_id"`
		TargetMessageID string `json:"target_message_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Bearer token from header
	token := ""
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}
	log.Printf("[Wodge Server] SecureChat Auth: HeaderLen=%d, TokenLen=%d", len(authHeader), len(token))

	stream, err := qastSvc.SecureChat(c.Request.Context(), req.Text, req.UserID, req.SessionID, req.TargetMessageID, token)
	if err != nil {
		log.Printf("[Wodge] SecureChat failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Flush immediately to establish stream
	c.Writer.Flush()

	// Manual copy loop to ensure flushing after every chunk
	buf := make([]byte, 1024)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			if _, wErr := c.Writer.Write(buf[:n]); wErr != nil {
				log.Printf("[Wodge] Streaming write error: %v", wErr)
				return // Client disconnected
			}
			c.Writer.Flush()
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("[Wodge] Streaming read error: %v", err)
			}
			break
		}
	}
}

// -- History Handlers --

func handleHistoryCreateSession(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
		Title  string `json:"title"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sess, err := qastSvc.CreateSession(c.Request.Context(), req.UserID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sess)
}

func handleHistoryGetSessions(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	sessions, err := qastSvc.GetSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func handleHistoryGetSession(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	sessionID := c.Param("id")
	sess, err := qastSvc.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		status := http.StatusInternalServerError
		// Naive check for 404
		if err.Error() == "failed to get session: 404" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sess)
}

func handleHistoryDeleteSession(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	sessionID := c.Param("id")
	if err := qastSvc.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// -- AstAuth Handlers --

// POST /api/auth/login
func handleAuthLogin(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := astAuthSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Sync User to QAST
	if qastSvc != nil {
		go func() {
			ctx := context.Background() // detach context
			if err := qastSvc.SyncUser(ctx, resp.User.ID, resp.User.Email, resp.User.Username, resp.User.FirstName, resp.User.LastName); err != nil {
				log.Printf("[Wodge] Failed to sync user %s to Qast: %v", resp.User.ID, err)
			}
		}()
	}

	c.JSON(http.StatusOK, resp)
}

// POST /api/auth/register
func handleAuthRegister(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}
	var req struct {
		Email           string `json:"email"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := astAuthSvc.Register(c.Request.Context(), req.Email, req.Username, req.Password, req.ConfirmPassword, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

// POST /api/auth/refresh
func handleAuthRefresh(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := astAuthSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/auth/verify
func handleAuthVerify(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}

	// Extract Bearer token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	// Strip "Bearer " prefix if present
	token := authHeader
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := astAuthSvc.VerifyToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Sync User to QAST (Async to not block response)
	if qastSvc != nil {
		go func() {
			ctx := context.Background()
			log.Printf("[Wodge] Syncing user %s (%s) to Qast...", user.ID, user.Username)
			if err := qastSvc.SyncUser(ctx, user.ID, user.Email, user.Username, user.FirstName, user.LastName); err != nil {
				log.Printf("[Wodge] Failed to sync user %s to Qast: %v", user.ID, err)
			} else {
				log.Printf("[Wodge] Successfully synced user %s to Qast", user.ID)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// POST /api/auth/logout
func handleAuthLogout(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := astAuthSvc.Logout(c.Request.Context(), req.AccessToken, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func handleHistoryShareSession(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	sessionID := c.Param("id")
	var req struct {
		TargetUsername string `json:"target_username"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := qastSvc.ShareSession(c.Request.Context(), sessionID, req.TargetUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func handleUsersSearch(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	resp, err := qastSvc.SearchUsers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func handleContextUpdate(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	id := c.Param("id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := qastSvc.UpdateContext(c.Request.Context(), id, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func handleContextGet(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}
	id := c.Param("id")
	ctxData, err := qastSvc.GetContext(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ctxData)
}
