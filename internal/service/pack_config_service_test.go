package service

import (
	"context"
	"errors"
	"testing"

	"go-packing/internal/domain"
)

type configRepoStub struct {
	cfg         *domain.PackConfig
	getErr      error
	saveErr     error
	savedCfg    domain.PackConfig
	savedPrev   int64
	saveInvoked bool
}

func (s *configRepoStub) Get(context.Context) (*domain.PackConfig, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.cfg == nil {
		return nil, nil
	}
	copyCfg := *s.cfg
	copyCfg.PackSizes = append([]int(nil), s.cfg.PackSizes...)
	return &copyCfg, nil
}

func (s *configRepoStub) SaveIfPreviousVersion(_ context.Context, cfg domain.PackConfig, previousVersion int64) error {
	s.saveInvoked = true
	s.savedCfg = cfg
	s.savedPrev = previousVersion
	if s.saveErr != nil {
		return s.saveErr
	}
	s.cfg = &cfg
	return nil
}

func TestPackConfigServiceGetCurrent(t *testing.T) {
	repo := &configRepoStub{}
	svc := NewPackConfigService(repo)

	cfg, err := svc.GetCurrent(context.Background())
	if err != nil {
		t.Fatalf("get current returned error: %v", err)
	}
	if cfg.Version != 0 {
		t.Fatalf("expected version 0, got %d", cfg.Version)
	}
	if len(cfg.PackSizes) != 0 {
		t.Fatalf("expected empty pack sizes, got %#v", cfg.PackSizes)
	}
}

func TestPackConfigServiceReplacePackSizes(t *testing.T) {
	repo := &configRepoStub{cfg: &domain.PackConfig{ID: 1, Version: 2, PackSizes: []int{10, 20}}}
	svc := NewPackConfigService(repo)

	cfg, err := svc.ReplacePackSizes(context.Background(), []int{53, 31, 23})
	if err != nil {
		t.Fatalf("replace returned error: %v", err)
	}

	if cfg.Version != 3 {
		t.Fatalf("expected version 3, got %d", cfg.Version)
	}
	if repo.savedPrev != 2 {
		t.Fatalf("expected previous version 2, got %d", repo.savedPrev)
	}
	if !repo.saveInvoked {
		t.Fatal("expected save to be invoked")
	}
	if len(cfg.PackSizes) != 3 || cfg.PackSizes[0] != 23 || cfg.PackSizes[2] != 53 {
		t.Fatalf("unexpected pack sizes: %#v", cfg.PackSizes)
	}
}

func TestPackConfigServiceReplacePackSizesConflict(t *testing.T) {
	repo := &configRepoStub{
		cfg:     &domain.PackConfig{ID: 1, Version: 1, PackSizes: []int{100}},
		saveErr: domain.ErrPackConfigVersionConflict,
	}
	svc := NewPackConfigService(repo)

	_, err := svc.ReplacePackSizes(context.Background(), []int{23, 31, 53})
	if !errors.Is(err, domain.ErrPackConfigVersionConflict) {
		t.Fatalf("expected version conflict, got %v", err)
	}
}

func TestPackConfigServiceReplacePackSizesValidation(t *testing.T) {
	repo := &configRepoStub{cfg: &domain.PackConfig{ID: 1, Version: 1, PackSizes: []int{100}}}
	svc := NewPackConfigService(repo)

	_, err := svc.ReplacePackSizes(context.Background(), []int{})
	if !errors.Is(err, domain.ErrInvalidPackSizes) {
		t.Fatalf("expected invalid pack sizes error, got %v", err)
	}
}
