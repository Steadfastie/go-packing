package domain

import "context"

type PackConfigRepository interface {
	Get(ctx context.Context) (*PackConfig, error)
	SaveIfPreviousVersion(ctx context.Context, cfg PackConfig, previousVersion int64) error
}
