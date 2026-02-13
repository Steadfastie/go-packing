package service

import (
	"reflect"
	"testing"

	"go-packing/internal/domain"
)

func TestOptimizePacks_TableExamples(t *testing.T) {
	packSizes := []int64{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name     string
		amount   int
		expected []domain.PackBreakdown
	}{
		{
			name:     "1 item",
			amount:   1,
			expected: []domain.PackBreakdown{{Size: 250, Count: 1}},
		},
		{
			name:     "250 items",
			amount:   250,
			expected: []domain.PackBreakdown{{Size: 250, Count: 1}},
		},
		{
			name:     "251 items",
			amount:   251,
			expected: []domain.PackBreakdown{{Size: 500, Count: 1}},
		},
		{
			name:   "501 items",
			amount: 501,
			expected: []domain.PackBreakdown{
				{Size: 500, Count: 1},
				{Size: 250, Count: 1},
			},
		},
		{
			name:   "12001 items",
			amount: 12001,
			expected: []domain.PackBreakdown{
				{Size: 5000, Count: 2},
				{Size: 2000, Count: 1},
				{Size: 250, Count: 1},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculate(tt.amount, packSizes)
			if err != nil {
				t.Fatalf("optimizePacks returned error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("unexpected result for amount=%d, got=%#v want=%#v", tt.amount, got, tt.expected)
			}
		})
	}
}
