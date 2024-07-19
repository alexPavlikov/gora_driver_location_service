package repository

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
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

func (r *Repo) ReadMessage(ctx context.Context, key int) error {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		return err
	}
	//defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition("drivers", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	//defer partConsumer.Close()

	var mu sync.Mutex
	var responseChannels = make(map[string]chan *sarama.ConsumerMessage)

	slog.Info("read message from kafka", "message key", key)

	var msg *sarama.ConsumerMessage

	go func() {

		for {
			var ok bool
			select {
			case msg, ok = <-partConsumer.Messages():
				fmt.Println(msg)
				if !ok {
					slog.Error("read from kafka channel closed, exiting goroutine")
					//	return
				}
				responseID := string(msg.Key)
				mu.Lock()
				ch, exists := responseChannels[responseID]
				if exists {
					ch <- msg
					delete(responseChannels, responseID)
				}
				mu.Unlock()
				fmt.Println(<-ch)
				fmt.Println(msg)
			}
		}

	}()

	slog.Info("read message from kafka", "message", msg)

	return nil
}
