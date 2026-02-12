package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go-packing/internal/domain/solver"
	"go-packing/internal/infrastructure/postgres"
	api "go-packing/internal/presentation/http"
	"go-packing/internal/service"
	"go-packing/pkg/logx"
)

func main() {
	logger := logx.NewJSONLogger()
	logger.Info("starting service")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/packing?sslmode=disable"
	}

	db, err := postgres.NewDB(context.Background(), dsn)
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
	optimizer := solver.NewOptimizer()

	calculateService := service.NewCalculateService(repo, optimizer)
	packConfigService := service.NewPackConfigService(repo)

	calculateHandler := api.NewCalculateHandler(calculateService, logger)
	packSizesHandler := api.NewPackSizesHandler(packConfigService, logger)

	router := api.NewRouter(logger, calculateHandler, packSizesHandler)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = fmt.Sprintf(":%s", port)
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
