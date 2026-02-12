package handlers

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

// NewPackSizesHandler builds handlers for /api/v1/pack-sizes endpoints.
func NewPackSizesHandler(svc *service.PackConfigService, logger *slog.Logger) *PackSizesHandler {
	return &PackSizesHandler{svc: svc, logger: logger}
}

// Get handles GET /api/v1/pack-sizes.
// @Summary Get current pack sizes
// @Description Returns configured pack sizes.
// @Tags Pack Sizes
// @Produce json
// @Success 200 {object} PackSizesResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/pack-sizes [get]
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

// Replace handles PUT /api/v1/pack-sizes.
// @Summary Replace pack sizes
// @Description Replaces all pack sizes and applies optimistic concurrency rules in persistence.
// @Tags Pack Sizes
// @Accept json
// @Produce json
// @Param request body PackSizesRequest true "Pack sizes payload"
// @Success 200 {object} PackSizesResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/pack-sizes [put]
func (h *PackSizesHandler) Replace(c *gin.Context) {
	var req PackSizesRequest

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
		// Conflict means another writer updated config between read and write.
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

	// Validation keeps persistence and solver assumptions simple.
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
