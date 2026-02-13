package domain

import (
	"sort"
	"time"
)

// PackConfig maps to the single persisted pack configuration row.
type PackConfig struct {
	Version   int64
	PackSizes []int64
	UpdatedAt time.Time
}

// NewPackConfig creates a new in-memory configuration.
func NewPackConfig(sizes []int64) (*PackConfig, error) {
	newSizes := make([]int64, len(sizes))
	copy(newSizes, sizes)
	sortPackSizesAsc(newSizes)

	return &PackConfig{
		Version:   0,
		PackSizes: newSizes,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

// Replace swaps pack sizes and advances version for CAS persistence updates.
func (p *PackConfig) Replace(sizes []int64) error {
	newSizes := make([]int64, len(sizes))
	copy(newSizes, sizes)
	sortPackSizesAsc(newSizes)

	p.PackSizes = newSizes
	p.Version++
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func sortPackSizesAsc(sizes []int64) {
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] < sizes[j]
	})
}
