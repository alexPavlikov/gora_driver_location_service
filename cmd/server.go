package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := setupLogger(cfg.Env)

	log.Info("starting application", "config", cfg)

	server.HandlerRequest()

	log.Info("initialization driver handlers")

	if err := http.ListenAndServe(cfg.ServerPath+":"+fmt.Sprint(cfg.ServerPort), nil); err != nil {
		log.Error("listen and serve server error", slog.Any("error", err.Error()))
		return err
	}

	return nil
}

// Функция для инициализации слоя логирования
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		slog.SetDefault(log)
	}

	return log
}
