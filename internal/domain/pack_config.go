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
	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	sort.Ints(newSizes)

	p.PackSizes = newSizes
	p.Version++
	p.UpdatedAt = time.Now().UTC()

	return nil
}