package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"go-packing/cmd/api/handlers"
	"go-packing/cmd/api/router"
	"go-packing/cmd/config"
	"go-packing/internal/infrastructure/postgres"
	"go-packing/internal/service"
	"go-packing/pkg/logx"
)

// @title Go Packing Service API
// @version 1.0
// @description API for calculating optimized pack allocations and managing pack-size configuration.
// @BasePath /
// @schemes http
func main() {
	cfg, err := config.Load()
	if err != nil {
		// Create a basic logger for errors before config is loaded
		slog.Error("config initialization failed", "error", err)
		os.Exit(1)
	}

	logger := logx.NewJSONLogger(cfg.Log.Level)
	logger.Info("starting service")
	logger.Info("configuration loaded", "env", cfg.AppEnv, "config_file", cfg.SourcePath)

	db, err := postgres.NewDB(context.Background(), cfg.Database.URL)
	if err != nil {
		logger.Error("database initialization failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("database close failed", "error", closeErr)
		}
	}()

	repo := postgres.NewPackConfigRepository(db, logger)

	calculateService := service.NewCalculateService(repo)
	packConfigService := service.NewPackConfigService(repo, logger)

	calculateHandler := handlers.NewCalculateHandler(calculateService, logger)
	packSizesHandler := handlers.NewPackSizesHandler(packConfigService, logger)

	router := router.NewRouter(logger, calculateHandler, packSizesHandler)

	addr := cfg.Server.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if gin.Mode() == gin.DebugMode {
		logger.Info("running in debug mode", "addr", addr)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}
