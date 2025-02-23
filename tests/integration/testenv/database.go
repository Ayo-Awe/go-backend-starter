package testenv

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ayo-awe/go-backend-starter/internal/db"
	_ "github.com/ayo-awe/go-backend-starter/internal/db/migrations"
	"github.com/ayo-awe/go-backend-starter/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB is the data access layer used for seeding data and interacting with
// the db without using the domain data access layer
type Database struct {
	// Existing app queries
	AppQueries *sqlc.Queries

	Pool *pgxpool.Pool
}

func (t *TestEnvironment) setupDatabase() error {
	pgContainer, err := postgres.Run(context.Background(),
		postgresImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return err
	}

	dbUrl, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	db, err := db.NewPostgresDatabase(context.Background(), dbUrl, t.Logger)
	if err != nil {
		return err
	}

	migrationsDir, err := filepath.Abs("../../../internal/db/migrations")
	if err != nil {
		return fmt.Errorf("failed to get migration path")
	}

	// run migrations
	if err := runMigrations(migrationsDir, dbUrl); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	t.postgres = pgContainer
	t.Database = db
	t.TestDB = &Database{AppQueries: sqlc.New(db.DB()), Pool: db.DB()}

	return nil
}

func runMigrations(migrationsDir string, databaseUrl string) error {
	db, err := goose.OpenDBWithDriver("pgx", databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to open db with driver: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
