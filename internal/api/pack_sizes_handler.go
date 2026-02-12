package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-packing/internal/domain"
	"go-packing/internal/service"
	"go-packing/pkg/httpx"
)

type PackSizesHandler struct {
	svc    *service.PackConfigService
	logger *slog.Logger
}

func NewPackSizesHandler(svc *service.PackConfigService, logger *slog.Logger) *PackSizesHandler {
	return &PackSizesHandler{svc: svc, logger: logger}
}

func (h *PackSizesHandler) Get(c *gin.Context) {
	packCfg, err := h.svc.GetCurrent(c.Request.Context())
	if err != nil {
		h.logger.Error("get pack sizes failed", "error", err)
		httpx.WriteError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	if packCfg == nil {
		c.JSON(http.StatusOK, gin.H{
			"pack_sizes": []int{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pack_sizes": packCfg.PackSizes,
	})
}

func (h *PackSizesHandler) Replace(c *gin.Context) {
	var req struct {
		PackSizes []int `json:"pack_sizes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}
	if !isValidPackSizes(req.PackSizes) {
		httpx.WriteError(c, http.StatusBadRequest, "INVALID_PACK_SIZES", domain.ErrInvalidPackSizes.Error())
		return
	}

	cfg, err := h.svc.ReplacePackSizes(c.Request.Context(), req.PackSizes)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConcurrencyConflict):
			httpx.WriteError(c, http.StatusConflict, "CONCURRENCY_CONFLICT", err.Error())
		default:
			h.logger.Error("replace pack sizes failed", "error", err)
			httpx.WriteError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pack_sizes": cfg.PackSizes,
	})
}

func isValidPackSizes(packSizes []int) bool {
	if len(packSizes) == 0 {
		return false
	}

	seen := make(map[int]struct{}, len(packSizes))
	for _, size := range packSizes {
		if size <= 0 {
			return false
		}
		if _, exists := seen[size]; exists {
			return false
		}
		seen[size] = struct{}{}
	}

	return true
}
