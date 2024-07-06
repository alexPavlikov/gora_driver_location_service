package locations

import (
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

type Handler struct {
	Producer sarama.SyncProducer
}

// DriverPostCord получающий координаты водителя
func (h *Handler) DriverPostCord(r *http.Request, data models.Cord) (map[string]string, error) {

	msg := sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder(fmt.Sprint(data.DriverID) + fmt.Sprint(data.Latitude) + fmt.Sprint(data.Longitude)),
	}

	if _, _, err := h.Producer.SendMessage(&msg); err != nil {
		return nil, fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return map[string]string{"status": "ok"}, nil
}
