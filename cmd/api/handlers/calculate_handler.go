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

type CalculateHandler struct {
	svc    *service.CalculateService
	logger *slog.Logger
}

// NewCalculateHandler builds a handler for POST /api/v1/calculate.
func NewCalculateHandler(svc *service.CalculateService, logger *slog.Logger) *CalculateHandler {
	return &CalculateHandler{svc: svc, logger: logger}
}

// Handle processes POST /api/v1/calculate and returns only the packs array.
// @Summary Calculate pack breakdown
// @Description Returns the optimal pack allocation for the requested amount.
// @Tags Calculate
// @Accept json
// @Produce json
// @Param request body CalculateRequest true "Calculation payload"
// @Success 200 {array} domain.PackBreakdown
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/calculate [post]
func (h *CalculateHandler) Handle(c *gin.Context) {
	var req CalculateRequest

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
		// Business rule: calculation requires configured pack sizes.
		case errors.Is(err, domain.ErrPackSizesNotConfigured):
			httpx.WriteError(c, http.StatusConflict, "PACK_SIZES_NOT_CONFIGURED", err.Error())
		case errors.Is(err, domain.ErrCouldNotCalculate):
			httpx.WriteError(c, http.StatusConflict, "COULD_NOT_CALCULATE", err.Error())
		default:
			h.logger.Error("calculate failed", "error", err)
			httpx.WriteError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, packs)
}
