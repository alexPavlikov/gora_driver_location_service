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

type wrappedFunc[Input, Output any] func(r *http.Request, data Input) (Output, error)

func handlerWrapper[Input, Output any](fn wrappedFunc[Input, Output]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Input
		//TODO decode

		response, err := fn(r, data)
		if err != nil {
			slog.ErrorContext(r.Context(), "can't handle request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_ = response
		//TODO encode response
	}
}
