package service

import (
	"context"
	"errors"
	"testing"

	"go-packing/internal/domain"
)

type calcRepoStub struct {
	cfg *domain.PackConfig
	err error
}

func (s calcRepoStub) Get(context.Context) (*domain.PackConfig, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.cfg == nil {
		return nil, nil
	}
	copyCfg := *s.cfg
	copyCfg.PackSizes = append([]int(nil), s.cfg.PackSizes...)
	return &copyCfg, nil
}

func (s calcRepoStub) Create(context.Context, domain.PackConfig) error {
	return nil
}

func (s calcRepoStub) Update(context.Context, domain.PackConfig) error {
	return nil
}

func TestCalculateService(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid amount", func(t *testing.T) {
		svc := NewCalculateService(calcRepoStub{})
		_, err := svc.Calculate(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}
	})

	t.Run("pack sizes not configured", func(t *testing.T) {
		svc := NewCalculateService(calcRepoStub{cfg: nil})
		_, err := svc.Calculate(ctx, 10)
		if !errors.Is(err, domain.ErrPackSizesNotConfigured) {
			t.Fatalf("expected ErrPackSizesNotConfigured, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := NewCalculateService(calcRepoStub{cfg: &domain.PackConfig{Version: 1, PackSizes: []int{250, 500, 1000}}})
		result, err := svc.Calculate(ctx, 251)
		if err != nil {
			t.Fatalf("calculate returned error: %v", err)
		}

		if len(result) != 1 || result[0].Size != 500 || result[0].Count != 1 {
			t.Fatalf("unexpected result: %#v", result)
		}
	})
}
