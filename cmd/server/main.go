package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ayo-awe/go-backend-starter/internal/app"
	"github.com/ayo-awe/go-backend-starter/internal/db"
	"github.com/ayo-awe/go-backend-starter/internal/server"
)

const service = "go-starter-api"

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger := setupLogger(cfg)
	ctx := context.Background()

	database, err := db.NewPostgresDatabase(ctx, cfg.DB.Dsn, logger)
	if err != nil {
		logger.Error("failed to setup database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	app := app.New()
	server := server.New(
		app,
		logger,
		server.WithCORS(cfg.Cors.AllowedOrigins),
		server.WithRateLimit(cfg.Limiter.Enabled, cfg.Limiter.Rate,
			cfg.Limiter.Window, cfg.Limiter.WhitelistedIPs,
		),
	)

	if err := server.Run(cfg.Port); err != nil {
		logger.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}

func setupLogger(cfg config) *slog.Logger {
	var handler slog.Handler

	output := os.Stdout

	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: cfg.Env == "development",
	}

	if cfg.Env == "development" {
		handler = slog.NewTextHandler(output, opts)
	} else {
		handler = slog.NewJSONHandler(output, opts)
	}

	return slog.New(handler).With(
		slog.String("service", service),
		slog.String("env", cfg.Env),
	)
}
