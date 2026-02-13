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

	result, err := calculate(amount, cfg.PackSizes)
	if err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	return result, nil
}

// calculate finds a pack combination that minimizes total shipped quantity,
// and among those, minimizes the number of packs.
func calculate(order int, packSizes []int64) ([]domain.PackBreakdown, error) {
	// Convert pack sizes to int and track the maximum size
	sizes := make([]int, len(packSizes))
	maxPack := 0
	for i, p := range packSizes {
		sizes[i] = int(p)
		maxPack = max(maxPack, sizes[i])
	}

	// Allow overfilling up to the largest pack size
	limit := order + maxPack

	// dp[i] = minimum number of packs needed to reach sum i
	// parent[i] = pack size last used to reach sum i
	dp := make([]int, limit+1)
	parent := make([]int, limit+1)

	for i := 1; i <= limit; i++ {
		dp[i] = math.MaxInt32
	}

	// Unbounded knapsack: minimize pack count for each achievable sum
	for _, s := range sizes {
		for i := s; i <= limit; i++ {
			// Rule #3: If using this pack reduces the total number of packs, update.
			if dp[i-s] != math.MaxInt32 && dp[i-s]+1 < dp[i] {
				dp[i] = dp[i-s] + 1
				parent[i] = s
			}
		}
	}

	// Rule #2: Choose the smallest reachable sum >= order (minimal overfill)
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

	// Time complexity: O((order + maxPack) * len(packSizes))
	// Space complexity: O(order + maxPack)
	return result, nil
}
