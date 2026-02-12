package domain

import (
	"sort"
	"time"
)

type PackConfig struct {
	Version   int64
	PackSizes []int
	UpdatedAt time.Time
}

func NewPackConfig(sizes []int) (*PackConfig, error) {
	if !isValidPackSizes(sizes) {
		return nil, ErrInvalidPackSizes
	}

	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	sort.Ints(newSizes)

	return &PackConfig{
		Version:   0,
		PackSizes: newSizes,
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (p *PackConfig) Replace(sizes []int) error {
	if !isValidPackSizes(sizes) {
		return ErrInvalidPackSizes
	}

	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	sort.Ints(newSizes)

	p.PackSizes = newSizes
	p.Version++
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func isValidPackSizes(sizes []int) bool {
	if len(sizes) == 0 {
		return false
	}

	seen := make(map[int]struct{}, len(sizes))
	for _, size := range sizes {
		if size <= 0 {
			return false
		}
		if _, exists := seen[size]; exists {
			return false
		}
		seen[size] = struct{}{}
	}

	return true
}
