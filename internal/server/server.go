package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayo-awe/go-backend-starter/internal/app"
	"github.com/ayo-awe/go-backend-starter/internal/oapi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

const (
	DefaultRateLimitEnabled = false
	DefaultRateLimitRate    = 100
	DefaultRateLimitWindow  = 1 * time.Minute
)

type Options func(s *Server)

func WithRateLimit(enabled bool, rate int, window time.Duration, whitelist []string) Options {
	return func(s *Server) {
		s.rateLimitEnabled = enabled
		s.rateLimitRate = rate
		s.rateLimitWindow = window
		s.rateLimitWhitelist = whitelist
	}
}

func WithCORS(allowedOrigins []string) Options {
	return func(s *Server) {
		s.corsAllowedOrigins = allowedOrigins
	}
}

type Server struct {
	logger             *slog.Logger
	app                *app.App
	rateLimitEnabled   bool
	rateLimitRate      int
	rateLimitWindow    time.Duration
	rateLimitWhitelist []string
	corsAllowedOrigins []string
}

func New(application *app.App, logger *slog.Logger, options ...Options) *Server {
	srv := &Server{
		logger:           logger,
		app:              application,
		rateLimitEnabled: DefaultRateLimitEnabled,
		rateLimitRate:    DefaultRateLimitRate,
		rateLimitWindow:  DefaultRateLimitWindow,
	}

	for _, opt := range options {
		if opt != nil {
			opt(srv)
		}
	}

	return srv
}

func (s Server) NewRouter() (http.Handler, error) {
	r := chi.NewRouter()

	requestLogger := httplog.Logger{
		Logger: s.logger,
		Options: httplog.Options{
			LogLevel:       slog.LevelDebug,
			JSON:           false,
			RequestHeaders: true,
		},
	}

	vm, err := s.ValidationMiddleware()
	if err != nil {
		return nil, err
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(httplog.Handler(&requestLogger))
	r.Use(s.CorsMiddleware())
	r.Use(s.RateLimitMiddleware())
	r.Use(vm)

	// maps routes to handlers using oapi spec
	h := oapi.HandlerFromMux(s, r)

	return h, nil
}

// Starts the server without blocking
func (s Server) Start(port int) (*http.Server, error) {
	r, err := s.NewRouter()
	if err != nil {
		return nil, fmt.Errorf("failed to create server router: %w", err)
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  r,
		ErrorLog: slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
	}

	go func() {
		s.logger.Info("Starting server...", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return srv, nil
}

// Starts the server and blocks until it receives a SIGINT or SIGTERM
func (s Server) Run(port int) error {
	srv, err := s.Start(port)
	if err != nil {
		return err
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownCh
	s.logger.Info("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("Graceful shutdown successful...")
	return nil
}
