package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/driver"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Функция инициализации и запуска сервера
func MustRun() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	driver.HandlerRequest()

	log.Info("initialization driver handlers")

	if err := http.ListenAndServe(cfg.ServerPath+":"+fmt.Sprint(cfg.ServerPort), nil); err != nil {
		log.Error("listen and serve server error", slog.Any("error", err.Error()))
		panic("listen and serve server error" + err.Error())
	}

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
	}

	return log
}
