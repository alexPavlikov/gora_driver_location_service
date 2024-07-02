package kafka

import (
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/alexPavlikov/gora_driver_location_service/internal/config"
)

// // код просто скопировал :)
// connect to kafka server
func Connect(cfg *config.Config) error {
	// responseChannels = make(map[string]chan *sarama.ConsumerMessage)

	// Создание продюсера Kafka
	producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Path + ":" + fmt.Sprint(cfg.Kafka.Port)}, nil)
	if err != nil {
		slog.Error("create producer kafka error" + err.Error())
		return err
	}
	defer producer.Close()

	// Создание консьюмера Kafka
	consumer, err := sarama.NewConsumer([]string{cfg.Kafka.Path + ":" + fmt.Sprint(cfg.Kafka.Port)}, nil)
	if err != nil {
		slog.Error("create consumer kafka error" + err.Error())
		return err
	}
	defer consumer.Close()

	// Подписка на партицию "pong" в Kafka
	partConsumer, err := consumer.ConsumePartition("pong", 0, sarama.OffsetNewest)
	if err != nil {
		slog.Error("subscribe to consume partition kafka error" + err.Error())
		return err
	}
	defer partConsumer.Close()

	return nil
}
