package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

type Repo struct {
	producer sarama.SyncProducer
	cfg      config.Config
}

func New(producer sarama.SyncProducer, cfg config.Config) *Repo {
	return &Repo{
		producer: producer,
		cfg:      cfg,
	}
}

func (r *Repo) SendMessage(ctx context.Context, cord models.CoordinatesPayload) error {
	var msg = sarama.ProducerMessage{
		Topic:     r.cfg.KafkaTopic,
		Key:       sarama.StringEncoder(fmt.Sprint(cord.Key)),
		Value:     sarama.StringEncoder(cord.Value),
		Timestamp: time.Now(),
	}

	slog.Info("send message to kafka", "message key", cord.Key)

	if _, _, err := r.producer.SendMessage(&msg); err != nil {
		return err
	}

	return nil
}
