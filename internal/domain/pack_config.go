package domain

import (
	"sort"
	"time"
)

type PackConfig struct {
	ID        int16
	Version   int64
	PackSizes []int
	UpdatedAt time.Time
}

func (p *PackConfig) Replace(newSizes []int) error {
	if !validPackSizes(newSizes) {
		return ErrInvalidPackSizes
	}

	next := make([]int, len(newSizes))
	copy(next, newSizes)
	sort.Ints(next)
	p.PackSizes = next
	p.Version++
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func validPackSizes(sizes []int) bool {
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
