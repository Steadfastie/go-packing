package service

import (
	"context"

	"go-packing/internal/domain"
)

type PackOptimizer interface {
	Optimize(amount int, packSizes []int) ([]domain.PackBreakdown, error)
}

type CalculateService struct {
	repo      domain.PackConfigRepository
	optimizer PackOptimizer
}

func NewCalculateService(repo domain.PackConfigRepository, optimizer PackOptimizer) *CalculateService {
	return &CalculateService{repo: repo, optimizer: optimizer}
}

func (s *CalculateService) Calculate(ctx context.Context, amount int) ([]domain.PackBreakdown, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	cfg, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	if cfg == nil || len(cfg.PackSizes) == 0 {
		return nil, domain.ErrPackSizesNotConfigured
	}

	return s.optimizer.Optimize(amount, cfg.PackSizes)
}
