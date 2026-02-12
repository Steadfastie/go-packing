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

type CalculateHandler struct {
	svc    *service.CalculateService
	logger *slog.Logger
}

func NewCalculateHandler(svc *service.CalculateService, logger *slog.Logger) *CalculateHandler {
	return &CalculateHandler{svc: svc, logger: logger}
}

func (h *CalculateHandler) Handle(c *gin.Context) {
	var req struct {
		Amount int `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}
	if req.Amount <= 0 {
		httpx.WriteError(c, http.StatusBadRequest, "INVALID_AMOUNT", domain.ErrInvalidAmount.Error())
		return
	}

	packs, err := h.svc.Calculate(c.Request.Context(), req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPackSizesNotConfigured):
			httpx.WriteError(c, http.StatusConflict, "PACK_SIZES_NOT_CONFIGURED", err.Error())
		default:
			h.logger.Error("calculate failed", "error", err)
			httpx.WriteError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, packs)
}
