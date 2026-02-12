package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-packing/internal/domain"
	"go-packing/internal/service"
)

type repoStub struct {
	getFn    func(context.Context) (*domain.PackConfig, error)
	createFn func(context.Context, domain.PackConfig) error
	updateFn func(context.Context, domain.PackConfig) error
}

func (s repoStub) Get(ctx context.Context) (*domain.PackConfig, error) {
	if s.getFn == nil {
		return nil, nil
	}
	return s.getFn(ctx)
}

func (s repoStub) Create(ctx context.Context, cfg domain.PackConfig) error {
	if s.createFn == nil {
		return nil
	}
	return s.createFn(ctx, cfg)
}

func (s repoStub) FindOneAndUpdate(ctx context.Context, cfg domain.PackConfig) error {
	if s.updateFn == nil {
		return nil
	}
	return s.updateFn(ctx, cfg)
}

func TestCalculateHandlerSuccess(t *testing.T) {
	repo := repoStub{getFn: func(context.Context) (*domain.PackConfig, error) {
		return &domain.PackConfig{Version: 1, PackSizes: []int{250, 500, 1000, 2000, 5000}}, nil
	}}
	r := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"amount":251}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var actual []domain.PackBreakdown
	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(actual) != 1 || actual[0].Size != 500 || actual[0].Count != 1 {
		t.Fatalf("unexpected response: %#v", actual)
	}
}

func TestCalculateHandlerNotConfigured(t *testing.T) {
	r := newTestRouter(repoStub{getFn: func(context.Context) (*domain.PackConfig, error) { return nil, nil }})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"amount":100}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestPackSizesHandlers(t *testing.T) {
	current := &domain.PackConfig{Version: 0, PackSizes: []int{}}
	repo := repoStub{
		getFn: func(context.Context) (*domain.PackConfig, error) {
			copyCfg := *current
			copyCfg.PackSizes = append([]int(nil), current.PackSizes...)
			return &copyCfg, nil
		},
		updateFn: func(_ context.Context, cfg domain.PackConfig) error {
			if cfg.Version-1 != current.Version {
				return domain.ErrConcurrencyConflict
			}
			current = &cfg
			return nil
		},
	}

	r := newTestRouter(repo)

	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, httptest.NewRequest(http.MethodGet, "/api/v1/pack-sizes", nil))
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wGet.Code)
	}

	wPut := httptest.NewRecorder()
	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/pack-sizes", bytes.NewBufferString(`{"pack_sizes":[23,31,53]}`))
	putReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wPut, putReq)
	if wPut.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", wPut.Code, wPut.Body.String())
	}
}

func TestPackSizesPutConflict(t *testing.T) {
	repo := repoStub{
		getFn: func(context.Context) (*domain.PackConfig, error) {
			return &domain.PackConfig{Version: 3, PackSizes: []int{100}}, nil
		},
		updateFn: func(context.Context, domain.PackConfig) error {
			return domain.ErrConcurrencyConflict
		},
	}
	r := newTestRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/pack-sizes", bytes.NewBufferString(`{"pack_sizes":[23,31,53]}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func newTestRouter(repo domain.PackConfigsRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	calcService := service.NewCalculateService(repo)
	packService := service.NewPackConfigService(repo, logger)

	calcHandler := NewCalculateHandler(calcService, logger)
	packHandler := NewPackSizesHandler(packService, logger)

	return NewRouter(logger, calcHandler, packHandler)
}
