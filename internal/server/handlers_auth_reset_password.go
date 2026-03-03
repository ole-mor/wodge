package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PUT /api/users/:id/reset-password
func handleAuthAdminResetPassword(c *gin.Context) {
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

	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := astAuthSvc.AdminResetPassword(c.Request.Context(), token, userID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
