package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DELETE /api/users/:id/delete-user
func handleAuthDeleteUser(c *gin.Context) {
	if astAuthSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AstAuth not configured"})
		return
	}

	userID := c.Param("id")
	token := ""
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	err := astAuthSvc.DeleteUser(c.Request.Context(), token, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
