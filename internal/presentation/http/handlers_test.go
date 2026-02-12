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
	"go-packing/internal/domain/solver"
	"go-packing/internal/service"
)

type repoStub struct {
	getFn  func(context.Context) (*domain.PackConfig, error)
	saveFn func(context.Context, domain.PackConfig, int64) error
}

func (s repoStub) Get(ctx context.Context) (*domain.PackConfig, error) {
	if s.getFn == nil {
		return nil, nil
	}
	return s.getFn(ctx)
}

func (s repoStub) SaveIfPreviousVersion(ctx context.Context, cfg domain.PackConfig, prev int64) error {
	if s.saveFn == nil {
		return nil
	}
	return s.saveFn(ctx, cfg, prev)
}

func TestCalculateHandlerSuccess(t *testing.T) {
	repo := repoStub{
		getFn: func(context.Context) (*domain.PackConfig, error) {
			return &domain.PackConfig{ID: 1, Version: 1, PackSizes: []int{250, 500, 1000, 2000, 5000}}, nil
		},
	}
	r := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"amount":251}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", w.Code, w.Body.String())
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
	current := &domain.PackConfig{ID: 1, Version: 0, PackSizes: []int{}}
	repo := repoStub{
		getFn: func(context.Context) (*domain.PackConfig, error) {
			copyCfg := *current
			copyCfg.PackSizes = append([]int(nil), current.PackSizes...)
			return &copyCfg, nil
		},
		saveFn: func(_ context.Context, cfg domain.PackConfig, prev int64) error {
			if prev != current.Version {
				return domain.ErrPackConfigVersionConflict
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

	var updated struct {
		Version   int   `json:"version"`
		PackSizes []int `json:"pack_sizes"`
	}
	if err := json.Unmarshal(wPut.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Version != 1 {
		t.Fatalf("expected version 1, got %d", updated.Version)
	}
}

func TestPackSizesPutConflict(t *testing.T) {
	repo := repoStub{
		getFn: func(context.Context) (*domain.PackConfig, error) {
			return &domain.PackConfig{ID: 1, Version: 3, PackSizes: []int{100}}, nil
		},
		saveFn: func(context.Context, domain.PackConfig, int64) error {
			return domain.ErrPackConfigVersionConflict
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

func newTestRouter(repo domain.PackConfigRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	calcService := service.NewCalculateService(repo, solver.NewOptimizer())
	packService := service.NewPackConfigService(repo)

	calcHandler := NewCalculateHandler(calcService, logger)
	packHandler := NewPackSizesHandler(packService, logger)

	return NewRouter(logger, calcHandler, packHandler)
}
