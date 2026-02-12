package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"go-packing/internal/domain"
)

type PackConfigRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPackConfigRepository(db *sql.DB, logger *slog.Logger) *PackConfigRepository {
	return &PackConfigRepository{db: db, logger: logger}
}

func (r *PackConfigRepository) Get(ctx context.Context) (*domain.PackConfig, error) {
	const query = `
SELECT id, version, pack_sizes, updated_at
FROM pack_configurations
WHERE id = 1
`

	row := r.db.QueryRowContext(ctx, query)

	var (
		cfg         domain.PackConfig
		packSizes32 []int32
	)
	if err := row.Scan(&cfg.ID, &cfg.Version, &packSizes32, &cfg.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("failed to fetch pack configuration", "error", err)
		return nil, fmt.Errorf("fetch pack configuration: %w", err)
	}

	cfg.PackSizes = convertToInt(packSizes32)
	return &cfg, nil
}

func (r *PackConfigRepository) SaveIfPreviousVersion(ctx context.Context, cfg domain.PackConfig, previousVersion int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const insertIfMissing = `
INSERT INTO pack_configurations (id, version, pack_sizes)
VALUES ($1, 0, '{}'::integer[])
ON CONFLICT (id) DO NOTHING
`

	if _, err := tx.ExecContext(ctx, insertIfMissing, cfg.ID); err != nil {
		r.logger.Error("failed to ensure pack configuration row", "error", err)
		return fmt.Errorf("ensure config row: %w", err)
	}

	const updateQuery = `
UPDATE pack_configurations
SET version = $1,
    pack_sizes = $2,
    updated_at = now()
WHERE id = $3
  AND version = $4
`

	result, err := tx.ExecContext(ctx, updateQuery, cfg.Version, convertToInt32(cfg.PackSizes), cfg.ID, previousVersion)
	if err != nil {
		r.logger.Error("failed to update pack configuration", "error", err)
		return fmt.Errorf("update pack configuration: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		// Optimistic concurrency: another writer moved the version ahead.
		return domain.ErrPackConfigVersionConflict
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func convertToInt(values []int32) []int {
	out := make([]int, len(values))
	for i := range values {
		out[i] = int(values[i])
	}
	return out
}

func convertToInt32(values []int) []int32 {
	out := make([]int32, len(values))
	for i := range values {
		out[i] = int32(values[i])
	}
	return out
}
