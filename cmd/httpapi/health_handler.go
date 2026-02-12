package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler serves a lightweight liveness endpoint.
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
