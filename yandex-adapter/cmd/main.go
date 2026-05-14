package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"yandex-adapter/config"
	"yandex-adapter/internal/app"
	"yandex-adapter/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(logger.Config{
		Level:  cfg.LogConfig.Level,
		Format: cfg.LogConfig.Format,
	})
	log.Info("starting application", "env", cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.New(cfg, log).Run(ctx); err != nil {
		log.Error("application terminated with error", "error", err)
		os.Exit(1)
	}
}
