package app

import (
	"context"
	"errors"
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
	yandexClient, err := a.buildYandexClient()
	if err != nil {
		return fmt.Errorf("init yandex client: %w", err)
	}

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

func (a *App) buildYandexClient() (geocodeuc.YandexClient, error) {
	if a.cfg.Yandex.UseStub {
		if a.cfg.Yandex.StubFile == "" {
			return nil, errors.New("YANDEX_STUB_FILE must be set when YANDEX_USE_STUB=true")
		}
		stub, err := yandex.NewStub()
		if err != nil {
			return nil, err
		}
		a.log.Warn("using yandex stub client")
		return stub, nil
	}

	if a.cfg.Yandex.APIKey == "" {
		return nil, errors.New("YANDEX_APIKEY must be set (or enable YANDEX_USE_STUB)")
	}
	return yandex.NewClient(yandex.Config{
		APIKey:  a.cfg.Yandex.APIKey,
		BaseURL: a.cfg.Yandex.BaseURL,
		Lang:    a.cfg.Yandex.Lang,
		Timeout: a.cfg.Yandex.Timeout,
	}), nil
}
