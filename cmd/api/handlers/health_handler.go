package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler serves a lightweight liveness endpoint.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
