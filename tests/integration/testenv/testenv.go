package testenv

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ayo-awe/go-backend-starter/internal/domain"
	"github.com/ayo-awe/go-backend-starter/internal/server"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbName        = "go-starter"
	dbPassword    = "go-starter"
	dbUser        = "go-starter"
	postgresImage = "postgres:16.3"
)

var (
	testenv *TestEnvironment = nil
)

type TestEnvironment struct {
	Database domain.Database
	TestDB   *Database
	Server   *server.Server
	Mailer   *mockMailer
	Logger   *slog.Logger
	postgres *postgres.PostgresContainer
}

func (t *TestEnvironment) Teardown() {
	// close db connection
	if t.TestDB != nil && t.TestDB.Pool != nil {
		t.TestDB.Pool.Close()
	}

	// cleanup postgres container
	if t.postgres != nil {
		if err := t.postgres.Terminate(context.Background()); err != nil {
			t.Logger.Error("failed to terminate postgres container",
				slog.String("error", err.Error()),
			)
		}
	}
}

// Setup should only be called once in TestMain. Use GetTestEnv is you're trying to access the
// test environment within a test or subtest
func Setup() (*TestEnvironment, error) {
	env := &TestEnvironment{}

	// setup mailer
	env.Mailer = newMockMailer()

	// setup logger
	env.Logger = slog.Default()

	// setup database
	if err := env.setupDatabase(); err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	// setup server
	if err := env.setupServer(); err != nil {
		return nil, fmt.Errorf("failed to setup server: %w", err)
	}

	testenv = env

	return env, nil
}

func Get(t *testing.T) *TestEnvironment {
	t.Cleanup(func() {
		testenv.Mailer.Clear()
	})

	return testenv
}
