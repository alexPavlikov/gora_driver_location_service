package server

import (
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type RouterBuilder struct {
	LocationsHandler *locations.Handler
}

func (r *RouterBuilder) Build() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/v1/locations", handlerWrapper(r.LocationsHandler.DriverPostCord))

	return router

}

type handlerWithError func(w http.ResponseWriter, r *http.Request) error

func handlerWrapper(fn handlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.ErrorContext(r.Context(), "can't handle request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
