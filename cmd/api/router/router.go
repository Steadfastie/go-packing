package router

import (
	"net/http"

	"log/slog"

	"go-packing/cmd/api/handlers"
	"go-packing/cmd/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-packing/docs"
)

// NewRouter wires HTTP routes, middleware, and Swagger UI.
func NewRouter(logger *slog.Logger, calculateHandler *handlers.CalculateHandler, packSizesHandler *handlers.PackSizesHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(logger))

	// Default entry point lands on Swagger UI.
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
	r.GET("/healthz", handlers.HealthHandler)
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Versioned API group for business endpoints.
	api := r.Group("/api/v1")
	api.POST("/calculate", calculateHandler.Handle)
	api.GET("/pack-sizes", packSizesHandler.Get)
	api.PUT("/pack-sizes", packSizesHandler.Replace)

	return r
}
