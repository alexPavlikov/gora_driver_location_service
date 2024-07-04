package kafka

import (
	"log/slog"

	"github.com/IBM/sarama"
)

// Функция получения producer
func GetProducer() (producer sarama.SyncProducer, err error) {
	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, nil) //....
	if err != nil {
		slog.Error("create producer kafka error - " + err.Error())
		return nil, err
	}

	return producer, nil
}

// Функция записи сообщения в kafka
func WriteMessage(producer sarama.SyncProducer, msg *sarama.ProducerMessage) error {
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		slog.Error("write message to kafka error - " + err.Error())
		return err
	}

	return nil
}
