package server

import (
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
)

// Инициализация всех handler сервера
func HandlerRequest() {
	http.HandleFunc("/driver_post_cord", locations.DriverPostCord)
}
