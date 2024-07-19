package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexPavlikov/gora_driver_location_service/internal/server/locations"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	DRIVER_ID string = "Driver_id"
)

type RouterBuilder struct {
	LocationsHandler *locations.Handler
}

func New(locationsHandler *locations.Handler) *RouterBuilder {
	return &RouterBuilder{
		LocationsHandler: locationsHandler,
	}
}

func (r *RouterBuilder) Build() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/v1/locations", middlewares(handlerWrapper(r.LocationsHandler.DriverPostCord)))
	router.Get("/v1/locations/{id}", mware(handlerWrapper(r.LocationsHandler.ReadDriverCordMessage)))

	return router
}

type wrappedFunc[Input, Output any] func(r *http.Request, data Input) (Output, error)

func handlerWrapper[Input, Output any](fn wrappedFunc[Input, Output]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Input

		if r.Method != http.MethodGet {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&data); err != nil {
				slog.ErrorContext(r.Context(), "can't decode data", "error", err)
				http.Error(w, "Bad request"+err.Error(), http.StatusBadRequest)
				return
			}
		}

		response, err := fn(r, data)
		if err != nil {
			slog.ErrorContext(r.Context(), "can't handle request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(&response); err != nil {
			slog.ErrorContext(r.Context(), "can't encode request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func mware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to convert driver_id", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), DRIVER_ID, id))
		h.ServeHTTP(w, r)
	}
}

func middlewares(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.Header.Get("X-ID"))
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to convert driver_id", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), DRIVER_ID, id))
		h.ServeHTTP(w, r)
	}
}
