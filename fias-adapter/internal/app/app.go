package app

import (
	"context"
	"fias-adapter/internal/service/fias"
	"fmt"
	"log/slog"

	"fias-adapter/config"
	httpctrl "fias-adapter/internal/controller/http"
	addressDetailUc "fias-adapter/internal/usecase/addressDetail"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{cfg: cfg, log: log}
}

func (a *App) Run(ctx context.Context) error {
	fiasClient, err := fias.NewClient(fias.Config{BaseURL: a.cfg.Fias.BaseURL, MasterToken: a.cfg.Fias.MasterToken})
	if err != nil {
		return fmt.Errorf("init yandex client: %w", err)
	}

	addressDetailUC := addressDetailUc.New(fiasClient)
	addressDetailHandler := httpctrl.NewAddressDetailHandler(addressDetailUC, a.log)

	router := httpctrl.NewRouter(addressDetailHandler)
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

//func (a *App) buildYandexClient() (addressDetailUc.FiasClient, error) {
//	if a.cfg.Yandex.UseStub {
//		stub, err := yandex.NewStub()
//		if err != nil {
//			return nil, err
//		}
//		a.log.Warn("using yandex stub client")
//		return stub, nil
//	}
//
//	return yandex.NewClient(yandex.Config{
//		APIKey:  a.cfg.Yandex.APIKey,
//		BaseURL: a.cfg.Yandex.BaseURL,
//		Lang:    a.cfg.Yandex.Lang,
//		Timeout: a.cfg.Yandex.Timeout,
//	})
//}
