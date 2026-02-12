package service

import (
	"context"

	"go-packing/internal/domain"
)

type PackConfigService struct {
	repo domain.PackConfigRepository
}

func NewPackConfigService(repo domain.PackConfigRepository) *PackConfigService {
	return &PackConfigService{repo: repo}
}

func (s *PackConfigService) GetCurrent(ctx context.Context) (domain.PackConfig, error) {
	cfg, err := s.repo.Get(ctx)
	if err != nil {
		return domain.PackConfig{}, err
	}
	if cfg == nil {
		return domain.PackConfig{ID: 1, Version: 0, PackSizes: []int{}}, nil
	}

	return *cfg, nil
}

func (s *PackConfigService) ReplacePackSizes(ctx context.Context, packSizes []int) (domain.PackConfig, error) {
	cfg, err := s.repo.Get(ctx)
	if err != nil {
		return domain.PackConfig{}, err
	}
	if cfg == nil {
		cfg = &domain.PackConfig{ID: 1, Version: 0, PackSizes: []int{}}
	}

	if err := cfg.Replace(packSizes); err != nil {
		return domain.PackConfig{}, err
	}

	previousVersion := cfg.Version - 1
	if err := s.repo.SaveIfPreviousVersion(ctx, *cfg, previousVersion); err != nil {
		return domain.PackConfig{}, err
	}

	return *cfg, nil
}
