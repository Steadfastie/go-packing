package solver

import (
	"testing"

	"go-packing/internal/domain"
)

func TestOptimizerChallengeExamples(t *testing.T) {
	op := NewOptimizer()
	sizes := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name     string
		amount   int
		expected map[int]int
	}{
		{name: "1", amount: 1, expected: map[int]int{250: 1}},
		{name: "250", amount: 250, expected: map[int]int{250: 1}},
		{name: "251", amount: 251, expected: map[int]int{500: 1}},
		{name: "501", amount: 501, expected: map[int]int{500: 1, 250: 1}},
		{name: "12001", amount: 12001, expected: map[int]int{5000: 2, 2000: 1, 250: 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := op.Optimize(tc.amount, sizes)
			if err != nil {
				t.Fatalf("optimize returned error: %v", err)
			}

			assertBreakdown(t, actual, tc.expected)
		})
	}
}

func TestOptimizerEdgeCaseLargeAmount(t *testing.T) {
	op := NewOptimizer()
	actual, err := op.Optimize(500000, []int{23, 31, 53})
	if err != nil {
		t.Fatalf("optimize returned error: %v", err)
	}

	expected := map[int]int{23: 2, 31: 7, 53: 9429}
	assertBreakdown(t, actual, expected)
}

func TestOptimizerInvalidInput(t *testing.T) {
	op := NewOptimizer()

	_, err := op.Optimize(0, []int{1})
	if err != domain.ErrInvalidAmount {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}

	_, err = op.Optimize(10, []int{})
	if err != domain.ErrInvalidPackSizes {
		t.Fatalf("expected ErrInvalidPackSizes, got %v", err)
	}
}

func assertBreakdown(t *testing.T, actual []domain.PackBreakdown, expected map[int]int) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("unexpected number of entries: got %d want %d", len(actual), len(expected))
	}

	for _, entry := range actual {
		want, ok := expected[entry.Size]
		if !ok {
			t.Fatalf("unexpected pack size in result: %d", entry.Size)
		}
		if entry.Count != want {
			t.Fatalf("unexpected count for size %d: got %d want %d", entry.Size, entry.Count, want)
		}
	}
}
