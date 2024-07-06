package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/kafka"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
)

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// setup logger
	config.SetupLogger(cfg.LogLevel)
	producer, err := kafka.GetProducer(cfg.Kafka.String())
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}

	// init handler request
	slog.Info("initialization driver handlers")
	locationsHandler := &locations.Handler{
		Producer: producer,
	}

	serverBuilder := server.RouterBuilder{
		LocationsHandler: locationsHandler,
	}
	srv := serverBuilder.Build()

	slog.Debug("starting application", "config", cfg)
	// load http server
	if err := http.ListenAndServe(cfg.Server.String(), srv); err != nil {
		slog.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}
