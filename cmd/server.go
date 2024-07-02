package server

import (
	"fmt"
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/kafka"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server"
)

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// setup logger
	log := config.SetupLogger(cfg.Env)

	log.Info("starting application", "config", cfg)

	// connect to kafka server
	err = kafka.Connect(cfg)
	if err != nil {
		log.Error("kafka connect error", "error", err.Error())
		return err
	}

	// init handler request
	server.HandlerRequest()

	log.Info("initialization driver handlers")

	// load http server
	if err := http.ListenAndServe(cfg.ServerPath+":"+fmt.Sprint(cfg.ServerPort), nil); err != nil {
		log.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}
