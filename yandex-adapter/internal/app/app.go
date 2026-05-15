package app

import (
	"context"
	"fmt"
	"log/slog"

	"yandex-adapter/config"
	httpctrl "yandex-adapter/internal/controller/http"
	"yandex-adapter/internal/service/yandex"
	geocodeuc "yandex-adapter/internal/usecase/geocode"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{cfg: cfg, log: log}
}

func (a *App) Run(ctx context.Context) error {
	yandexClient := yandex.NewClient(yandex.Config{
		APIKey:  a.cfg.Yandex.APIKey,
		BaseURL: a.cfg.Yandex.BaseURL,
		Lang:    a.cfg.Yandex.Lang,
		Timeout: a.cfg.Yandex.Timeout,
	})

	geocodeUC := geocodeuc.New(yandexClient)
	geocodeHandler := httpctrl.NewGeocodeHandler(geocodeUC, a.log)

	router := httpctrl.NewRouter(geocodeHandler)
	server := httpctrl.NewServer(a.cfg.HttpServer, router)

	serverErr := make(chan error, 1)
	go func() {
		a.log.Info("http server started", "addr", server.Addr())
		serverErr <- server.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	case <-ctx.Done():
		a.log.Info("shutdown signal received, stopping http server")
	}

	shutdownCtx := context.Background()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	a.log.Info("http server stopped gracefully")
	return nil
}
