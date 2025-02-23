package testenv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type Seeder interface {
	Seed(ctx context.Context, db *Database) error
}

type SeederFunc func(ctx context.Context, db *Database) error

func (s SeederFunc) Seed(ctx context.Context, db *Database) error {
	return s(ctx, db)
}

func (t *TestEnvironment) SeedDatabase(tt *testing.T, seeder Seeder) {
	ctx := context.Background()
	require.NoError(tt, seeder.Seed(ctx, t.TestDB))

	tt.Cleanup(func() {
		// Clears all tables and restarts the serial sequences
		require.NoError(tt, t.TestDB.AppQueries.TruncateAllTables(ctx))
	})
}
