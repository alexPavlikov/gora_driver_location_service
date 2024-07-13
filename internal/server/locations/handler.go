package locations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
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

type emptyResponse struct{}

// Handler получающий координаты водителя
func (h *Handler) DriverPostCord(r *http.Request, data models.Cord) (emptyResponse, error) {

	ctx := r.Context()

	topic := ctx.Value("Topic")

	driverID := ctx.Value("Driver_id")

	data.DriverID = driverID.(int)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return emptyResponse{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	msg := sarama.ProducerMessage{
		Topic: topic.(string),
		Key:   sarama.StringEncoder(data.DriverID),
		Value: sarama.StringEncoder(jsonData),
	}

	if err = h.Service.SendMessage(&msg); err != nil {
		return emptyResponse{}, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return emptyResponse{}, nil
}
