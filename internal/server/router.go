package server

import (
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
	"github.com/go-chi/chi"
)

// Инициализация всех handler сервера
func HandlerRequest(router *chi.Mux) {

	router.Post("/driver_post_cord", locations.DriverPostCord)
}
