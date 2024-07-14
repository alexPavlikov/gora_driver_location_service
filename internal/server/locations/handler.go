package locations

import (
	"encoding/json"
	"fmt"
	"net/http"

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

type emptyResponse struct{}

// Handler получающий координаты водителя
func (h *Handler) DriverPostCord(r *http.Request, data DriverPostCordRequest) (emptyResponse, error) {

	ctx := r.Context()

	_, err := json.Marshal(data)
	if err != nil {
		return emptyResponse{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	var msg = models.Cord{
		Longitude: data.Longitude,
		Latitude:  data.Latitude,
	}

	if err = h.Service.StoreMessage(ctx, msg); err != nil {
		return emptyResponse{}, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return emptyResponse{}, nil
}
