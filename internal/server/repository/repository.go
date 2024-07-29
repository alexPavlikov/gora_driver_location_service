package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
)

type Repo struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	cfg      config.Config
}

func New(cfg config.Config, producer sarama.SyncProducer, consumer sarama.Consumer) *Repo {
	return &Repo{
		producer: producer,
		consumer: consumer,
		cfg:      cfg,
	}
}

func (r *Repo) SendMessage(ctx context.Context, cord models.CoordinatesPayload) error {

	cordJSON, err := json.Marshal(cord)
	if err != nil {
		return fmt.Errorf("failed to marshal cord: %w", err)
	}

	var msg = sarama.ProducerMessage{
		Topic:     r.cfg.KafkaTopic,
		Key:       sarama.StringEncoder(fmt.Sprint(cord.Key)),
		Value:     sarama.ByteEncoder(cordJSON),
		Timestamp: time.Now(),
	}

	slog.Info("send message to kafka", "message key", cord.Key)

	if _, _, err := r.producer.SendMessage(&msg); err != nil {
		return err
	}

	return nil
}

func (r *Repo) ReadMessageFromKafka() ([]models.Cord, error) {

	partition, err := r.consumer.Partitions(r.cfg.KafkaTopic)
	if err != nil {
		return nil, err
	}

	partitionConsumer, err := r.consumer.ConsumePartition(r.cfg.KafkaTopic, partition[0], sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}
	defer partitionConsumer.Close()

	var cords []models.Cord

	for c := range partitionConsumer.Messages() {
		log.Printf("Consumed message: [%s], offset: [%d]\n", c.Value, c.Offset)
		var cord models.Cord
		if err := json.Unmarshal(c.Value, &cord); err != nil {
			return nil, fmt.Errorf("failed unmarshal result from kafka: %w", err)
		}
		cords = append(cords, cord)
	}

	return cords, nil
}
