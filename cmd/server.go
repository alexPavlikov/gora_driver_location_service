package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/kafka"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/service"
)

// input

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// setup logger
	config.SetupLogger(cfg.LogLevel)

	slog.Info("starting application", "config", cfg)

	producer, err := kafka.GetProducer(cfg.Kafka.ToString())
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}

	defer producer.Close()

	// init handler request
	slog.Info("initialization driver handlers")

	locationsService := &service.Service{
		Producer: producer,
	}

	locationsHandler := &locations.Handler{
		Service: locationsService,
	}

	serverBuilder := server.RouterBuilder{
		LocationsHandler: locationsHandler,
	}

	srv := serverBuilder.Build()

	// load http server
	if err := http.ListenAndServe(cfg.Server.ToString(), srv); err != nil {
		slog.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}
