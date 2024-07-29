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

	srv, close, err := NewServer(cfg)
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

func NewServer(cfg *config.Config) (http.Handler, func() error, error) {

	slog.Info("starting application", "server config", cfg.Server.ToString())

	// get producer for kafka
	producer, err := kafka.GetProducer(cfg.Kafka.ToString())
	if err != nil {
		return nil, nil, fmt.Errorf("error creating producer: %w", err)
	}

	consumer, close, err := kafka.GetConsumer(cfg.Kafka.ToString())
	if err != nil {
		return nil, nil, fmt.Errorf("failed run server - kafka get consumer: %w", err)
	}

	defer close()

	// init handler request
	slog.Info("initialization driver handlers")

	repository := repository.New(*cfg, producer, consumer)

	service := service.New(repository)

	handlers := locations.New(service)

	serverBuilder := server.New(handlers)

	srv := serverBuilder.Build()

	return srv, producer.Close, nil
}
