package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lib/pq"

	"go-packing/internal/domain"
)

type PackConfigRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewPackConfigRepository creates a PostgreSQL-backed configuration repository.
func NewPackConfigRepository(db *sql.DB, logger *slog.Logger) *PackConfigRepository {
	return &PackConfigRepository{db: db, logger: logger}
}

// Get loads the single pack config row (id=1). It returns nil when not initialized.
func (r *PackConfigRepository) Get(ctx context.Context) (*domain.PackConfig, error) {
	const query = `
		SELECT version, COALESCE(pack_sizes, '{}'::INTEGER[]), updated_at
		FROM pack_configs
		WHERE id = 1
	`

	row := r.db.QueryRowContext(ctx, query)

	var packCfg domain.PackConfig
	if err := row.Scan(&packCfg.Version, pq.Array(&packCfg.PackSizes), &packCfg.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("failed to fetch pack config", "error", err)
		return nil, fmt.Errorf("fetch pack config: %w", err)
	}

	return &packCfg, nil
}

// Create inserts the initial config row if it does not already exist.
func (r *PackConfigRepository) Create(ctx context.Context, packCfg domain.PackConfig) error {
	const insertQuery = `
		INSERT INTO pack_configs (id, pack_sizes, version, updated_at)
		VALUES (1, $1, $2, $3)
		ON CONFLICT DO NOTHING
	`

	_, err := r.db.ExecContext(
		ctx,
		insertQuery,
		pq.Array(packCfg.PackSizes),
		packCfg.Version,
		packCfg.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("failed to create pack config", "error", err)
		return fmt.Errorf("create pack config: %w", err)
	}

	return nil
}

// Update performs an optimistic-concurrency update using version CAS.
func (r *PackConfigRepository) Update(ctx context.Context, packCfg domain.PackConfig) error {
	const updateQuery = `
		UPDATE pack_configs
		SET pack_sizes = $1,
			version = $2,
			updated_at = $3
		WHERE id = 1
			AND version = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		updateQuery,
		pq.Array(packCfg.PackSizes),
		packCfg.Version,
		packCfg.UpdatedAt,
		packCfg.Version-1,
	)
	if err != nil {
		r.logger.Error("failed to update pack config", "error", err)
		return fmt.Errorf("update pack config: %w", err)
	}

	// No affected rows means version mismatch (concurrent writer won).
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrConcurrencyConflict
	}

	return nil
}
