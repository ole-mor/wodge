package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PUT /api/users/:id/activate-user
func handleAuthActivateUser(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}

	userID := c.Param("id")

	// Extract Bearer token from header
	token := ""
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	var req struct {
		Active bool `json:"active"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := astAuthSvc.ActivateUser(c.Request.Context(), token, userID, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
