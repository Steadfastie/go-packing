package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"go-packing/internal/api"
	"go-packing/internal/config"
	"go-packing/internal/infrastructure/postgres"
	"go-packing/internal/service"
	"go-packing/pkg/logx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// Create a basic logger for errors before config is loaded
		slog.Error("config initialization failed", "error", err)
		os.Exit(1)
	}
	
	logger := logx.NewJSONLogger(cfg.Log.Level)
	logger.Info("starting service")

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

	calculateHandler := api.NewCalculateHandler(calculateService, logger)
	packSizesHandler := api.NewPackSizesHandler(packConfigService, logger)

	router := api.NewRouter(logger, calculateHandler, packSizesHandler)

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
