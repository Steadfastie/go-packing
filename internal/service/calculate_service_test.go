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
	return s.cfg, nil
}

func (s calcRepoStub) SaveIfPreviousVersion(context.Context, domain.PackConfig, int64) error {
	return nil
}

type optimizerStub struct {
	result []domain.PackBreakdown
	err    error
}

func (o optimizerStub) Optimize(int, []int) ([]domain.PackBreakdown, error) {
	if o.err != nil {
		return nil, o.err
	}
	return o.result, nil
}

func TestCalculateService(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid amount", func(t *testing.T) {
		svc := NewCalculateService(calcRepoStub{}, optimizerStub{})
		_, err := svc.Calculate(ctx, 0)
		if !errors.Is(err, domain.ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}
	})

	t.Run("pack sizes not configured", func(t *testing.T) {
		svc := NewCalculateService(calcRepoStub{cfg: nil}, optimizerStub{})
		_, err := svc.Calculate(ctx, 10)
		if !errors.Is(err, domain.ErrPackSizesNotConfigured) {
			t.Fatalf("expected ErrPackSizesNotConfigured, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		expected := []domain.PackBreakdown{{Size: 500, Count: 1}}
		svc := NewCalculateService(
			calcRepoStub{cfg: &domain.PackConfig{ID: 1, Version: 1, PackSizes: []int{250, 500}}},
			optimizerStub{result: expected},
		)

		actual, err := svc.Calculate(ctx, 251)
		if err != nil {
			t.Fatalf("calculate returned error: %v", err)
		}
		if len(actual) != 1 || actual[0].Size != 500 || actual[0].Count != 1 {
			t.Fatalf("unexpected result: %#v", actual)
		}
	})
}
