package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"go-packing/internal/domain"
)

type packConfigRepoStub struct {
	current *domain.PackConfig
	getErr  error
	addErr  error
	updErr  error

	created *domain.PackConfig
	updated *domain.PackConfig
}

func (s *packConfigRepoStub) Get(context.Context) (*domain.PackConfig, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.current == nil {
		return nil, nil
	}
	copyCfg := *s.current
	copyCfg.PackSizes = append([]int(nil), s.current.PackSizes...)
	return &copyCfg, nil
}

func (s *packConfigRepoStub) Create(_ context.Context, cfg domain.PackConfig) error {
	s.created = &cfg
	if s.addErr != nil {
		return s.addErr
	}
	s.current = &cfg
	return nil
}

func (s *packConfigRepoStub) FindOneAndUpdate(_ context.Context, cfg domain.PackConfig) error {
	s.updated = &cfg
	if s.updErr != nil {
		return s.updErr
	}
	s.current = &cfg
	return nil
}

func newPackConfigServiceForTest(repo domain.PackConfigsRepository) *PackConfigService {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return NewPackConfigService(repo, logger)
}

func TestPackConfigServiceCreate(t *testing.T) {
	repo := &packConfigRepoStub{}
	svc := newPackConfigServiceForTest(repo)

	cfg, err := svc.ReplacePackSizes(context.Background(), []int{53, 23, 31})
	if err != nil {
		t.Fatalf("replace returned error: %v", err)
	}

	if cfg.Version != 0 {
		t.Fatalf("expected version 0, got %d", cfg.Version)
	}
	if repo.created == nil {
		t.Fatal("expected create to be called")
	}
}

func TestPackConfigServiceUpdate(t *testing.T) {
	repo := &packConfigRepoStub{
		current: &domain.PackConfig{Version: 2, PackSizes: []int{250, 500}},
	}
	svc := newPackConfigServiceForTest(repo)

	cfg, err := svc.ReplacePackSizes(context.Background(), []int{1000, 500})
	if err != nil {
		t.Fatalf("replace returned error: %v", err)
	}

	if cfg.Version != 3 {
		t.Fatalf("expected version 3, got %d", cfg.Version)
	}
	if repo.updated == nil {
		t.Fatal("expected update to be called")
	}
}

func TestPackConfigServiceConflict(t *testing.T) {
	repo := &packConfigRepoStub{
		current: &domain.PackConfig{Version: 1, PackSizes: []int{250}},
		updErr:  domain.ErrConcurrencyConflict,
	}
	svc := newPackConfigServiceForTest(repo)

	_, err := svc.ReplacePackSizes(context.Background(), []int{500})
	if !errors.Is(err, domain.ErrConcurrencyConflict) {
		t.Fatalf("expected ErrConcurrencyConflict, got %v", err)
	}
}
