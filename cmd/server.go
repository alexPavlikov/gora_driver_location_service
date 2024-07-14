package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/kafka"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/repository"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/service"
)

// input

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	srv, close, err := ServerLoad(cfg)
	if err != nil {
		return err
	}

	defer close()

	// load http server
	if err := http.ListenAndServe(cfg.Server.ToString(), srv); err != nil {
		slog.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}

func ServerLoad(cfg *config.Config) (http.Handler, func() error, error) {
	// setup logger
	config.SetupLogger(cfg.LogLevel)

	slog.Info("starting application", "server config", cfg.Server.ToString())

	// get producer for kafka
	producer, err := kafka.GetProducer(cfg.Kafka.ToString())
	if err != nil {
		return nil, nil, fmt.Errorf("error creating producer: %w", err)
	}

	// init handler request
	slog.Info("initialization driver handlers")

	repository := repository.New(producer, *cfg)

	service := service.New(repository)

	handlers := locations.New(service)

	serverBuilder := server.New(handlers)

	srv := serverBuilder.Build()

	return srv, producer.Close, nil
}
