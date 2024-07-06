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

		id, err := strconv.Atoi(r.Header.Get("Driver_id"))
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to convert header value", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		topic := r.Header.Get("Topic")

		ctx := context.WithValue(context.Background(), "Driver_id", id)
		ctx = context.WithValue(ctx, "Topic", topic)

		r = r.WithContext(ctx)

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			slog.ErrorContext(r.Context(), "can't decode data", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response, err := fn(r, data)
		if err != nil {
			slog.ErrorContext(r.Context(), "can't handle request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		_ = response
		// encode response
	}
}
