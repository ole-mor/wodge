package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	
	"test-app2/internal/handlers"
)

func main() {
	r := gin.Default()
	
	// Register generated routes
	handlers.RegisterRoutes(r)

	log.Println("Starting backend on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
