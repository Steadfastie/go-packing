package service

import (
	"context"
	"log/slog"

	"go-packing/internal/domain"
)

type PackConfigService struct {
	repo   domain.PackConfigsRepository
	logger *slog.Logger
}

func NewPackConfigService(repo domain.PackConfigsRepository, logger *slog.Logger) *PackConfigService {
	return &PackConfigService{repo: repo, logger: logger}
}

func (s *PackConfigService) GetCurrent(ctx context.Context) (*domain.PackConfig, error) {
	cfg, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *PackConfigService) ReplacePackSizes(ctx context.Context, packSizes []int) (*domain.PackConfig, error) {
	packCfg, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	if packCfg == nil {
		s.logger.Info("pack config not found, creating new one")

		packCfg, err = domain.NewPackConfig(packSizes)
		if err != nil {
			return nil, err
		}

		if err := s.repo.Create(ctx, *packCfg); err != nil {
			return nil, err
		}
		return packCfg, nil
	}

	s.logger.Info("pack config found, updating existing one", "version", packCfg.Version)
	if err := packCfg.Replace(packSizes); err != nil {
		return nil, err
	}

	if err := s.repo.FindOneAndUpdate(ctx, *packCfg); err != nil {
		return nil, err
	}

	return packCfg, nil
}
