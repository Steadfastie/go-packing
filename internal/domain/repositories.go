package domain

import "context"

type PackConfigsRepository interface {
	Get(ctx context.Context) (*PackConfig, error)
	Create(ctx context.Context, packCfg PackConfig) error
	FindOneAndUpdate(ctx context.Context, packCfg PackConfig) error
}
