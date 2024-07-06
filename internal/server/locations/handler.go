package locations

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

type Handler struct {
	Producer sarama.SyncProducer
}

// DriverPostCord получающий координаты водителя
func (h *Handler) DriverPostCord(w http.ResponseWriter, r *http.Request) error {
	var cord models.Cord
	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cord); err != nil {
		slog.Error("driver post cord error decode cord - " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	msg := sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder(fmt.Sprint(cord.DriverID) + fmt.Sprint(cord.Latitude) + fmt.Sprint(cord.Longitude)),
	}

	if _, _, err := h.Producer.SendMessage(&msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}
