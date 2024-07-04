package locations

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/kafka"
	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

// Handler получающий координаты водителя
func DriverPostCord(w http.ResponseWriter, r *http.Request) {
	slog.Info("request DriverPostCord")

	var cord models.Cord

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		slog.Error("driver post cord error content type")
		w.WriteHeader(http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cord); err != nil {
		slog.Error("driver post cord error decode cord - " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	msg := sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder(fmt.Sprint(cord.DriverID) + fmt.Sprint(cord.Latitude) + fmt.Sprint(cord.Longitude)),
	}

	producer, err := kafka.GetProducer()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	defer producer.Close()

	if err = kafka.WriteMessage(producer, &msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	slog.Info("write message completed")
}
