package main

import (
	"log/slog"
	"os"

	"github.com/tuanbt/go-be-template/internal/app"
	"github.com/tuanbt/go-be-template/internal/config"
	"github.com/tuanbt/go-be-template/internal/platform/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := log.New(cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(logger)

	if err := app.RunWithConfig(cfg, logger); err != nil {
		logger.Error("fatal", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
