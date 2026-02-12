package domain

import "testing"

func TestPackConfigReplace(t *testing.T) {
	cfg := PackConfig{
		ID:        1,
		Version:   3,
		PackSizes: []int{100, 200},
	}

	err := cfg.Replace([]int{53, 23, 31})
	if err != nil {
		t.Fatalf("replace returned error: %v", err)
	}

	if cfg.Version != 4 {
		t.Fatalf("expected version 4, got %d", cfg.Version)
	}

	expected := []int{23, 31, 53}
	for i := range expected {
		if cfg.PackSizes[i] != expected[i] {
			t.Fatalf("unexpected pack sizes at %d: got %d want %d", i, cfg.PackSizes[i], expected[i])
		}
	}
}

func TestPackConfigReplaceValidation(t *testing.T) {
	tests := []struct {
		name  string
		sizes []int
	}{
		{name: "empty", sizes: []int{}},
		{name: "negative", sizes: []int{10, -1}},
		{name: "zero", sizes: []int{10, 0}},
		{name: "duplicate", sizes: []int{10, 10}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := PackConfig{}
			err := cfg.Replace(tc.sizes)
			if err == nil {
				t.Fatalf("expected validation error")
			}
			if err != ErrInvalidPackSizes {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
