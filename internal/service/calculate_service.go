package service

import (
	"context"
	"math"
	"sort"

	"go-packing/internal/domain"
)

type CalculateService struct {
	repo domain.PackConfigsRepository
}

func NewCalculateService(repo domain.PackConfigsRepository) *CalculateService {
	return &CalculateService{repo: repo}
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

	return optimizePacks(amount, cfg.PackSizes)
}

func optimizePacks(amount int, packSizes []int) ([]domain.PackBreakdown, error) {
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	minPack := sizes[len(sizes)-1]
	upper := amount + minPack - 1
	inf := math.MaxInt / 4

	dp := make([]int, upper+1)
	prevTotal := make([]int, upper+1)
	prevPack := make([]int, upper+1)

	for i := range dp {
		dp[i] = inf
		prevTotal[i] = -1
		prevPack[i] = -1
	}
	dp[0] = 0

	for total := 1; total <= upper; total++ {
		for _, pack := range sizes {
			if total-pack < 0 {
				continue
			}
			if dp[total-pack] == inf {
				continue
			}

			candidate := dp[total-pack] + 1
			if candidate < dp[total] {
				dp[total] = candidate
				prevTotal[total] = total - pack
				prevPack[total] = pack
			}
		}
	}

	bestTotal := -1
	for total := amount; total <= upper; total++ {
		if dp[total] != inf {
			bestTotal = total
			break
		}
	}
	if bestTotal == -1 {
		return nil, domain.ErrPackSizesNotConfigured
	}

	counts := make(map[int]int, len(sizes))
	for cur := bestTotal; cur > 0; {
		pack := prevPack[cur]
		if pack <= 0 {
			return nil, domain.ErrInvalidPackSizes
		}
		counts[pack]++
		cur = prevTotal[cur]
	}

	result := make([]domain.PackBreakdown, 0, len(counts))
	for _, pack := range sizes {
		if counts[pack] > 0 {
			result = append(result, domain.PackBreakdown{Size: pack, Count: counts[pack]})
		}
	}

	return result, nil
}
