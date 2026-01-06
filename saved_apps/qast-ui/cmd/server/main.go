package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"qast-ui/internal/handlers"
)

func main() {
	r := gin.Default()

	// Register generated routes
	handlers.RegisterRoutes(r)

	log.Println("Starting backend on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
