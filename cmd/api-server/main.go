package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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

	// Add your API routes here
	// Example: r.GET("/api/users", getUsersHandler)

	log.Println("Starting API server on :8080")
	log.Println("Frontend will access APIs via: http://localhost:5173/api")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
