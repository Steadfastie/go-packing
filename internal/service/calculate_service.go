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

// NewCalculateService creates a calculation service backed by pack configuration storage.
func NewCalculateService(repo domain.PackConfigsRepository) *CalculateService {
	return &CalculateService{repo: repo}
}

// Calculate returns an optimal pack breakdown for the requested amount.
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

	return calculate(amount, cfg.PackSizes)
}

// calculate minimizes shipped quantity first, then pack count.
func calculate(order int, packSizes []int64) ([]domain.PackBreakdown, error) {
	// Find maximum pack size to bound DP range
	sizes := make([]int, len(packSizes))
	maxPack := 0
	for i, p := range packSizes {
		sizes[i] = int(p)
		maxPack = max(maxPack, sizes[i])
	}

	limit := order + maxPack

	// dp[i] = min packs to reach sum i
	// parent[i] = the size of the last pack used to reach sum i
	dp := make([]int, limit+1)
	parent := make([]int, limit+1)

	for i := 1; i <= limit; i++ {
		dp[i] = math.MaxInt32
	}

	// Unbounded knapsack DP. Complexity: O(len(packSizes) * limit)
	for _, s := range sizes {
		for i := s; i <= limit; i++ {
			// Rule #3: If using this pack results in fewer total packs, update.
			if dp[i-s] != math.MaxInt32 && dp[i-s]+1 < dp[i] {
				dp[i] = dp[i-s] + 1
				parent[i] = s
			}
		}
	}

	// Rule #2: find smallest overfill, then fewest packs
	bestSum := -1
	for i := order; i <= limit; i++ {
		if dp[i] != math.MaxInt32 {
			bestSum = i
			break
		}
	}

	if bestSum == -1 {
		return nil, domain.ErrCouldNotCalculate
	}

	// Reconstruct pack selection
	counts := make(map[int]int)
	curr := bestSum
	for curr > 0 {
		p := parent[curr]
		counts[p]++
		curr -= p
	}

	result := make([]domain.PackBreakdown, 0, len(counts))
	for size, count := range counts {
		result = append(result, domain.PackBreakdown{Size: size, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	// Time Complexity: O((order + maxPack) * number_of_pack_sizes)
	// Space Complexity: O(order + maxPack)
	return result, nil
}
