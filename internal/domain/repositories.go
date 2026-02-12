package domain

import "context"

// PackConfigsRepository persists and retrieves the single pack configuration document.
type PackConfigsRepository interface {
	Get(ctx context.Context) (*PackConfig, error)
	Create(ctx context.Context, packCfg PackConfig) error
	FindOneAndUpdate(ctx context.Context, packCfg PackConfig) error
}
