package domain

import "testing"

func TestNewPackConfig(t *testing.T) {
	cfg, err := NewPackConfig([]int{53, 23, 31})
	if err != nil {
		t.Fatalf("new pack config returned error: %v", err)
	}

	if cfg.Version != 0 {
		t.Fatalf("expected version 0, got %d", cfg.Version)
	}
	if len(cfg.PackSizes) != 3 || cfg.PackSizes[0] != 23 || cfg.PackSizes[2] != 53 {
		t.Fatalf("unexpected pack sizes: %#v", cfg.PackSizes)
	}
}

func TestPackConfigReplace(t *testing.T) {
	cfg, err := NewPackConfig([]int{250, 500})
	if err != nil {
		t.Fatalf("new pack config returned error: %v", err)
	}

	if err := cfg.Replace([]int{1000, 500, 250}); err != nil {
		t.Fatalf("replace returned error: %v", err)
	}

	if cfg.Version != 1 {
		t.Fatalf("expected version 1, got %d", cfg.Version)
	}
	if len(cfg.PackSizes) != 3 || cfg.PackSizes[0] != 250 || cfg.PackSizes[2] != 1000 {
		t.Fatalf("unexpected pack sizes: %#v", cfg.PackSizes)
	}
}

func TestPackConfigValidation(t *testing.T) {
	if _, err := NewPackConfig([]int{}); err != ErrInvalidPackSizes {
		t.Fatalf("expected ErrInvalidPackSizes, got %v", err)
	}

	cfg, err := NewPackConfig([]int{1})
	if err != nil {
		t.Fatalf("new pack config returned error: %v", err)
	}
	if err := cfg.Replace([]int{1, 1}); err != ErrInvalidPackSizes {
		t.Fatalf("expected ErrInvalidPackSizes, got %v", err)
	}
}
