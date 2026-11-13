package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanbt/go-be-template/internal/config"
	httpserver "github.com/tuanbt/go-be-template/internal/http"
	"github.com/tuanbt/go-be-template/internal/platform/health"
	"github.com/tuanbt/go-be-template/internal/platform/postgres"
	"github.com/tuanbt/go-be-template/internal/store/db"
	"github.com/tuanbt/go-be-template/internal/transport"
	"github.com/tuanbt/go-be-template/internal/user"
)

// App is the composition root for the HTTP API process.
type App struct {
	log   *slog.Logger
	pool  *pgxpool.Pool
	srv   *http.Server
	ready *health.Gate
}

// New builds dependencies without starting the HTTP server.
func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*App, http.Handler, error) {
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres: %w", err)
	}

	queries := db.New(pool)
	repo := user.NewSQLRepository(queries)
	userSvc := user.NewService(repo)
	userHTTP := user.NewHandler(log, userSvc)

	ready := health.NewGate(postgres.NewPoolChecker(pool))
	api := transport.NewAPI(log, ready, userHTTP)
	router := httpserver.NewRouter(log, api)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: httpserver.ReadHeaderTimeout,
		ReadTimeout:       httpserver.ReadTimeout,
		WriteTimeout:      httpserver.WriteTimeout,
		IdleTimeout:       httpserver.IdleTimeout,
	}

	return &App{log: log, pool: pool, srv: srv, ready: ready}, router, nil
}

// RunWithConfig starts the API and blocks until graceful shutdown.
func RunWithConfig(cfg config.Config, logger *slog.Logger) error {
	ctx := context.Background()
	application, _, err := New(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer application.pool.Close()
	return application.serve(logger)
}

// HandlerForTest returns a wired handler and pool (caller runs migrations, then closes pool).
func HandlerForTest(ctx context.Context, cfg config.Config, log *slog.Logger) (http.Handler, *pgxpool.Pool, error) {
	application, handler, err := New(ctx, cfg, log)
	if err != nil {
		return nil, nil, err
	}
	return handler, application.pool, nil
}

func (a *App) serve(logger *slog.Logger) error {
	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", slog.String("addr", a.srv.Addr))
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		logger.Info("shutdown", slog.String("signal", sig.String()))
	}

	// Fail readiness before draining so a load balancer can stop sending new work.
	a.ready.Close()
	drain := time.NewTimer(httpserver.ReadinessDrain)
	defer drain.Stop()
	select {
	case <-drain.C:
	case sig := <-sigCh:
		logger.Info("shutdown forced", slog.String("signal", sig.String()))
		return a.srv.Close()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpserver.ShutdownTimeout)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- a.srv.Shutdown(shutdownCtx)
	}()
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		return nil
	case sig := <-sigCh:
		logger.Info("shutdown forced", slog.String("signal", sig.String()))
		return a.srv.Close()
	}
}
