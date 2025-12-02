// internal/handlers/health.go (ตัวอย่าง)
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/impk123/Inventory-Management-Mini-System/pkg/database"
)

func HealthCheck(c *gin.Context) {

	if err := database.HealthCheck(database.DB); err != nil {

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "Database unhealthy: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "System is healthy",
	})
}
