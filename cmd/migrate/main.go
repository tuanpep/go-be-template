package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tuanbt/go-be-template/internal/config"
	"github.com/tuanbt/go-be-template/internal/platform/migrate"
	"github.com/tuanbt/go-be-template/internal/platform/postgres"
	"github.com/tuanbt/go-be-template/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := postgres.Ping(pingCtx, pool); err != nil {
		slog.Error("postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := migrate.Up(ctx, pool, migrations.FS, "."); err != nil {
		slog.Error("migrate", slog.String("error", err.Error()))
		os.Exit(1)
	}
	fmt.Println("migrations applied")
}
