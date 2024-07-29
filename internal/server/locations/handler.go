package locations

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/service"
)

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

type DriverPostCordRequest struct {
	Longitude float32 `json:"longitude"` // ширина
	Latitude  float32 `json:"latitude"`  // долгота
}

type DriverPostCordResponse struct {
	ID        int     `json:"id"`        // id водителя
	Longitude float32 `json:"longitude"` // ширина
	Latitude  float32 `json:"latitude"`  // долгота
}

type emptyResponse struct{}

// Handler получающий координаты водителя
func (h *Handler) DriverPostCord(r *http.Request, data DriverPostCordRequest) (emptyResponse, error) {

	ctx := r.Context()

	var msg = models.Cord{
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	}

	if err := h.Service.StoreMessage(ctx, msg); err != nil {
		return emptyResponse{}, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return emptyResponse{}, nil
}

func (h *Handler) ReadDriverCordMessage(r *http.Request, data emptyResponse) ([]DriverPostCordResponse, error) {
	ctx := r.Context()

	var cordsRep = make([]DriverPostCordResponse, 0)

	var mu sync.Mutex

	var ch = make(chan models.Cord)

	go func() {

		mu.Lock()

		cords, err := h.Service.ReadMessage(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "POST MessageAddHandler", "error", err)
		}

		for _, v := range cords {
			ch <- v
		}

		mu.Unlock()
	}()

	go func() {

		for v := range ch {
			mu.Lock()

			var msg = DriverPostCordResponse{
				ID:        v.DriverID,
				Longitude: v.Longitude,
				Latitude:  v.Latitude,
			}

			cordsRep = append(cordsRep, msg)
			mu.Unlock()
		}

	}()

	fmt.Println(cordsRep)

	return cordsRep, nil
}
