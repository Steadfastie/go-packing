package domain

import (
	"sort"
	"time"
)

// PackConfig maps to the single persisted pack configuration row.
type PackConfig struct {
	Version   int64
	PackSizes []int
	UpdatedAt time.Time
}

// NewPackConfig creates a new in-memory configuration.
func NewPackConfig(sizes []int) (*PackConfig, error) {
	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	sort.Ints(newSizes)

	return &PackConfig{
		Version:   0,
		PackSizes: newSizes,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

// Replace swaps pack sizes and advances version for CAS persistence updates.
func (p *PackConfig) Replace(sizes []int) error {
	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	sort.Ints(newSizes)

	p.PackSizes = newSizes
	p.Version++
	p.UpdatedAt = time.Now().UTC()

	return nil
}
