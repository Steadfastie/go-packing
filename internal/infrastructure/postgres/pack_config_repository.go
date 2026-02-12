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

func NewPackConfigRepository(db *sql.DB, logger *slog.Logger) *PackConfigRepository {
	return &PackConfigRepository{db: db, logger: logger}
}

func (r *PackConfigRepository) Get(ctx context.Context) (*domain.PackConfig, error) {
	const query = `
		SELECT version, pack_sizes, updated_at
		FROM pack_configs
		WHERE id = 1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query)

	var cfg domain.PackConfig
	if err := row.Scan(&cfg.Version, pq.Array(&cfg.PackSizes), &cfg.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("failed to fetch pack configuration", "error", err)
		return nil, fmt.Errorf("fetch pack configuration: %w", err)
	}

	return &cfg, nil
}

func (r *PackConfigRepository) Create(ctx context.Context, cfg domain.PackConfig) error {
	const insertQuery = `
		INSERT INTO pack_configs (id, pack_sizes, version, updated_at)
		VALUES (1, $1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, insertQuery, pq.Array(cfg.PackSizes), cfg.Version, cfg.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConcurrencyConflict
		}
		r.logger.Error("failed to create pack configuration", "error", err)
		return fmt.Errorf("create pack configuration: %w", err)
	}

	return nil
}

func (r *PackConfigRepository) FindOneAndUpdate(ctx context.Context, cfg domain.PackConfig) error {
	const updateQuery = `
		UPDATE pack_configs
		SET pack_sizes = $1,
			version = $2,
			updated_at = $3
		WHERE id = 1
			AND version = $4
	`

	previousVersion := cfg.Version - 1
	result, err := r.db.ExecContext(ctx, updateQuery, pq.Array(cfg.PackSizes), cfg.Version, cfg.UpdatedAt, previousVersion)
	if err != nil {
		r.logger.Error("failed to update pack configuration", "error", err)
		return fmt.Errorf("update pack configuration: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrConcurrencyConflict
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}

	return pqErr.Code == "23505"
}
