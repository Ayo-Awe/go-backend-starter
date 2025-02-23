package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ayo-awe/go-backend-starter/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDatabase struct {
	dbPool  *pgxpool.Pool
	queries *sqlc.Queries
	logger  *slog.Logger
}

func (p *PostgresDatabase) DB() *pgxpool.Pool {
	return p.dbPool
}

func NewPostgresDatabase(ctx context.Context, dsn string, logger *slog.Logger) (*PostgresDatabase, error) {
	// create connection pool
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	db := &PostgresDatabase{
		dbPool:  pool,
		queries: queries,
		logger:  logger,
	}

	return db, nil
}

func (p *PostgresDatabase) Close() {
	p.dbPool.Close()
}

func (p *PostgresDatabase) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := p.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin database transaction: %w", err)
	}

	defer func(ctx context.Context, tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			p.logger.Error("failed to rollback database transaction", "err", err)
		}
	}(ctx, tx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit database transactoin: %w", err)
	}
	return nil
}
