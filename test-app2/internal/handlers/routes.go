package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes is called by main to setup routes
// This file will be updated by wodge generator
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
