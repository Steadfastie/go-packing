package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func NewRouter(logger *slog.Logger, calculateHandler *CalculateHandler, packSizesHandler *PackSizesHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(logger))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.POST("/calculate", calculateHandler.Handle)
	api.GET("/pack-sizes", packSizesHandler.Get)
	api.PUT("/pack-sizes", packSizesHandler.Replace)

	return r
}
