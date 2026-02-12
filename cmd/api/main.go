package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	api "go-packing/internal/presentation/http"
	"go-packing/pkg/logx"
)

func main() {
	logger := logx.NewJSONLogger()
	logger.Info("starting service")

	r := api.NewRouter(logger)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}
