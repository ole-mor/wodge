package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PUT /api/history/experiment/state
func handleHistoryUpdateUserExperimentState(c *gin.Context) {
	if qastSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "QAST not configured"})
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		State  string `json:"state"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Wodge] Updating experiment state for user %s to %s", req.UserID, req.State)
	result, err := qastSvc.UpdateUserExperimentState(c.Request.Context(), req.UserID, req.State)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
